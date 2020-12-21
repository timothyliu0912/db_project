$("#search").click(function () {
    var query = $("#query").val();
    var db_type = $('#db_select').val();
    $.ajax({
        type: "GET",
        url: "/search",
        data: {
            query: query,
            type: db_type
        },
        success: function (msg) {
            let urls = [];
            let publish = [];
            $("#videos div#video").remove();
            console.log(msg.data);
            msg.data.forEach(item => {
                const add = '<div class="card text-center" id="video">'+
                    '<div class="embed-responsive embed-responsive-16by9">'+
                    '<iframe class="embed-responsive-item" src="'+item.vid+'"></iframe>'+
                    '</div>'+
                    '<div class="card-body">'+
                    '<h6 class="card-title">'+ item.title +'</h6>'+
                    '<p class="card-publish text-sm-left font-weight-lighter">時間：'+item.published+'</p>'+
                    '<p class="card-viewcount text-sm-left font-weight-lighter">觀看次數：'+ item.viewcount +'</p>'+
                    '<p class="card-author text-sm-left font-weight-lighter">'+ item.Author +'</p>'+
                    '<a href="'+item.url+'" class="btn btn-outline-primary stretched-link">View</a>'+
                    '</div>'+
                    '</div>';
                $("#videos").append(add);
            })
        }

    });

});