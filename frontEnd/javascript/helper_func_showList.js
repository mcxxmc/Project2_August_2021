let urlShowList = "http://localhost:8080/showList/";

function appendRows(data) {
    let t = $(".content #div_records #table_myTable").DataTable();
    for (let i = 0; i < data["records"].length; ++ i){
        let record = data["records"][i];
        t.row.add([
            record["name"],
            record["path"],
            record["prediction"],
            record["label"]
        ]).node().id = record["id"];
        t.draw(false);
    }
}

$(window).on("load", function(){
    $.ajax({
        type: "GET",
        url: urlShowList,
        success:function(data){
            appendRows(data);
        },
        error: function(data){
            alert("error!");
            console.log("error");
            console.log(data);
        }
    });
});