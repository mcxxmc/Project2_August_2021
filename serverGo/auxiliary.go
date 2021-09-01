package main

import "fmt"

// checkErr checks the error and prints it out if not nil
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}