let urlPostImage = 'http://localhost:8080/imgSystem/'  // must specify http to enable CORS

$(document).ready(function (e) {
    $('#form_img').on('submit',(function(e) {
        e.preventDefault();
        var formData = new FormData(this);

        $.ajax({
            type:'POST',
            url: urlPostImage,
            data:formData,
            cache:false,
            contentType: false,
            processData: false,
            success:function(data){
                alert(data["r"])
                console.log("success");
                console.log(data);
            },
            error: function(data){
                alert("error!")
                console.log("error");
                console.log(data);
            }
        });
    }));
});
