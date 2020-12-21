$(":button.special-btn").click(function () {
    var Tbdata = {}; 
    var td_length = $(this).parents("tr").children().length; 
    var id = $(this).parents("tr").children('td:eq(0)').data('id');
    console.log(id);
    $.ajax({
        type: "GET",
        url: "/show_all",
        data: {
            query: id
        },
         success: function (msg) {
            $("#service_title").val(msg.data.title); 
            console.log(msg.data);
         }

     });
    
    // $('#edit').append(add);
});

  