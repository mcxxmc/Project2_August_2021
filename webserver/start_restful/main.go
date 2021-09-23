package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"webserver/common"
	"webserver/db"
	"webserver/webservice"
)

func addMiddleware(router *gin.Engine) {
	router.Use(webservice.Filter())
	router.Use(webservice.SetCORS())
}

func bindUrl(router *gin.Engine) {
	router.POST("/upload", webservice.PostImage)
	router.POST("/prediction", webservice.ImmediatePred)
	router.POST("/pictures", webservice.ShowPictures)
	router.POST("/labels-pictures", webservice.PostImageLabels)

	router.GET("/list", webservice.ShowList)
	router.GET("/labels-pictures", webservice.GetUnlabeledPictures)
	// TODO: safety / authorization
	router.GET("/opencv", webservice.UseCamera)
}

func createServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr: addr,
		Handler: handler,
	}
}

func run(server *http.Server) {
	zap.S().Info("Server started, listening at: ", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.S().Fatalf("listen: %s\n", err)
	}
}

func gracefulShutDown(server *http.Server, delay time.Duration) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<- quit  // block here
	zap.S().Infof("Server is closing...")
	ctx, cancel := context.WithTimeout(context.Background(), delay)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server shutdown: ", err)
	}
	zap.S().Infof("Server exits.")
}

func main() {
	// initiate the zap logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()  // flushes buffer, if any
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// create the shared connection pool to the database
	db.OpenSharedDb()
	defer db.CloseSharedDb()

	// Test the connection to the database and check if the table exists.
	db.TryConnection()

	// Set up the router.
	// By default, for Go > 1.8, http server will use all available CPUs.
	// If it is going to be run in a linux container, use "go.uber.org/automaxprocs" to catch the max number of CPUs
	// in the virtual machine.
	numCPUs := runtime.NumCPU()
	zap.S().Infof("Golang running, number of logical CPUs usable: %d\n", numCPUs)
	// runtime.GOMAXPROCS(4)
	router := gin.Default()
	addMiddleware(router)
	bindUrl(router)

	// config the server
	server := createServer(common.WebserverPortGin, router)

	// Run the server at port 8080 (by default).
	go run(server)

	// gracefully shut down (5 seconds delay)
	gracefulShutDown(server, 5 * time.Second)
}
