$(document).ready(function(){
    M.AutoInit();

    function reset() {
        $("input.reset").val('');
    }

    $('.modal').modal();

    $("#copy_public_key").click(function() {
        var key = $("#self_public_key");
        key.select();
        document.execCommand("copy");
        M.toast({html: 'Public Key Copied', classes: 'green'});
    })

    $("#add_user_form").submit(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/addUser',
            data: $('form#add_user_form').serialize(),
            success: function(data) {
                M.toast({html: data, classes: 'green'});
                var instance = M.Modal.getInstance($("#add_user_modal"));
                instance.close();
                reset();
            },
            error: function(error) {
                M.toast({html: error.responseText, classes: 'red'});
            }
        })
    })
});