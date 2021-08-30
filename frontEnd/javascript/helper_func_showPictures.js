let urlShowPictures = "http://localhost:8080/showPictures/";
let lastPictureId = 0;  // the id of the last picture
let defaultN = 5;


function ajax_showPictures(offset, n){
    $.ajax({
        type: 'POST',
        url: urlShowPictures,
        data: JSON.stringify({"offset": offset, "n": n}),  
        // should not specify things like content-type here for the backend is golang-gin
        success: function(data){
            let imagebundles = data["images"];
            if (imagebundles == null) {
                alert("No more pictures available.");
            }else {
                let div = $("#div_pictures");
                for (let i = 0; i < imagebundles.length; i ++){
                    let imagebundle = imagebundles[i];
                    div.append('<img src="data:image/png;base64, ' + imagebundle["image"] + 
                    '" /><br><p>' + imagebundle["text"] + '</p>');
                }
            }
        },
        error: function(data){
            alert("error!");
            console.log("error");
            console.log(data);
        }
    });
};

function displayMore(n=defaultN){
    ajax_showPictures(lastPictureId, n);
    lastPictureId += n;
};

$(window).on("load", function() {
    displayMore();
});

$(function () {
    $("#button_next").on("click", function() {
    displayMore();
    });
})
