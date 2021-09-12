let urlLabelPictures = "http://localhost:8080/label-pictures";
let allPictureNames = [];

// display the pictures on the web page with check boxes
function displayPicturesAndOptions(imagebundles){
    if (imagebundles == null){
        alert("No pictures for labeling.");
    }else{
        let div = $("#div_pictures");
        for (let i = 0; i < imagebundles.length; i ++){
            let imagebundle = imagebundles[i];
            let name = imagebundle["text"];
            allPictureNames.push(name);
            div.append("<img src='data:image/png;base64, " + imagebundle["image"] + "' /><br>");
            div.append("<p>" + name + "</p><br>");
            div.append("<fieldSet><legend>You Choice:</legend>" + 
            "<input type='radio' name='radio_" + name + "' value='v'>Vehicle<br>" + 
            "<input type='radio' name='radio_" + name + "' value='nv'>Non-vehicle<br>" + 
            "</fieldSet><br>");
        }
    }
};

// get the pictures and display them on the web page
function ajax_requestPictures_GET() {
    $.ajax({
        type: "GET",
        url: urlLabelPictures,
        success: function(data){
            let imagebundles = data["images"];
            displayPicturesAndOptions(imagebundles);
        },
        error: function(data){
            alert("error!");
            console.log("error");
            console.log(data);
        }
    });
} ;

// construct a json object from the users response
function collectRadioChoices(){
    let r = {"results": []};
    let name = "";
    let val = "";
    for (let i = 0; i < allPictureNames.length; i ++){
        name = allPictureNames[i];
        val = $('input[name="radio_' + name + '"]:checked').val();
        // check if the val is empty
        if (val != null){
            let tmp = {"name": name, "val": val};
            r["results"].push(tmp);
        }
    }
    alert(JSON.stringify(r))
    return r;
}

// upload the labeled results
function ajax_uploadLabeledResults_POST() {
    if (allPictureNames == null) {
        alert("No new data to send.");
    }else{
        let results = collectRadioChoices();
        $.ajax({
            type: "POST",
            url: urlLabelPictures,
            data: JSON.stringify(results),
            success: function(data) {
                alert("Results successfully sent!");
            },
            error: function(data) {
                alert("error!");
                console.log("error");
                console.log(data);
            }
        })
    }
}

// initialization
$(window).on("load", function() {
    ajax_requestPictures_GET();
});

$(function () {

    // on click button
    $("#button_submit").on("click", function() {
        ajax_uploadLabeledResults_POST();
    });
});

