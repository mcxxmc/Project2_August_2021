let urlOpenCV = "http://localhost:8080/opencv";


// make the ajax call
function make_ajax_openCV() {
    $.ajax({
        type: "GET",
        url: urlOpenCV,
        success: function(data){
            alert("Success!")
        },
        error: function(data){
            alert("error!");
            console.log("error");
            console.log(data);
        }
    });
}

// initialization
$(window).on("load", function() {
    make_ajax_openCV();
})
