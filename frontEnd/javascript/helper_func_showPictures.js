let urlShowPictures = "http://localhost:8080/showPictures/";
let lastPictureId = 0;  // the id of the last picture
let defaultN = 5;

// display the pictures and the texts on the web page
function displayPictures(imagebundles) {
    let div = $("#div_pictures");
    for (let i = 0; i < imagebundles.length; i++) {
        let imagebundle = imagebundles[i];
        div.append('<img src="data:image/png;base64, ' + imagebundle["image"] +
            '" /><br><p>' + imagebundle["text"] + '</p>');
    }
}

// Ask for new photos (by passing to backend the id of the latest received photo)
function ajax_showPictures(offset, n) {
    $.ajax({
        type: 'POST',
        url: urlShowPictures,
        data: JSON.stringify({ "offset": offset, "n": n }),
        // should not specify things like content-type here for the backend is golang-gin
        success: function (data) {
            let imagebundles = data["images"];
            if (imagebundles == null) {
                alert("No more pictures available.");
            } else {
                displayPictures(imagebundles);
            }
        },
        error: function (data) {
            alert("error!");
            console.log("error");
            console.log(data);
        }
    });
};

// automatically manage / update the ids
function displayMore(n = defaultN) {
    ajax_showPictures(lastPictureId, n);
    lastPictureId += n;
};

// initialization
$(window).on("load", function () {
    displayMore();
});

$(function () {

    // on click button
    $("#button_next").on("click", function () {
        displayMore();
    });
})
