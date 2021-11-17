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

    $("#add_server_form").submit(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/addServer',
            data: $('form#add_server_form').serialize(),
            success: function(data) {
                M.toast({html: data, classes: 'green'});
                var instance = M.Modal.getInstance($("#add_server_modal"));
                instance.close();
                reset();
            },
            error: function(error) {
                M.toast({html: error.responseText, classes: 'red'});
            }
        })
    })

    $("#list_users").click(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/getUsers',
            success: function(data) {
                var users = JSON.parse(data)
                $("#list_users_tbody").html('');
                users.forEach(function (user) {
                    $("#list_users_tbody").append(
                        '<tr>\
                        <td>'+user.id+'</td>\
                        <td>'+user.username+'</td>\
                        <td>'+user.email+'</td>\
                        <td>'+user.public_key+'</td>\
                        <td>'+user.created_at+'</td>\
                        <td><button class="btn waves-effect waves-light red delete_user_button" type="submit" data-userid="'+user.id+'">Delete<i class="material-icons right">delete</i></button></td>\
                        </tr>')
                });
                
            },
            error: function(error) {
                M.toast({html: error.responseText, classes: 'red'});
            }
        })
    })

    $("#list_servers").click(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/getServers',
            success: function(data) {
                var servers = JSON.parse(data)
                $("#list_servers_tbody").html('');
                servers.forEach(function (server) {
                    $("#list_servers_tbody").append('<tr>\
                    <td>'+server.id+'</td>\
                    <td>'+server.username+'</td>\
                    <td>'+server.ip+'</td>\
                    <td>'+server.created_at+'</td>\
                    <td><button class="btn waves-effect waves-light red delete_server_button" type="submit" data-serverid="'+server.id+'">Delete<i class="material-icons right">delete</i></button></td>\
                    </tr>')
                });
                
            },
            error: function(error) {
                M.toast({html: error.responseText, classes: 'red'});
            }
        })
    })    
    
    
    $("#list_access").click(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/getAccess',
            success: function(data) {
                console.log(data)
                var accesslist = JSON.parse(data)

                var userAccessArray = {};
                if(accesslist.accesses) {
                    accesslist.accesses.forEach(function (access) {
                        if(!userAccessArray[access.user_id]) {
                            userAccessArray[access.user_id] = []
                        }
                        (userAccessArray[access.user_id]).push(access.server_id);
                    });
                }
                
                $("#list_access_headings").html('<th></th>');
                accesslist.servers.forEach(function (server) {
                        $("#list_access_headings").append('<th>'+server.username+'</th>')
                });

                $("#list_access_tbody").html('');
                accesslist.users.forEach(function (user) {
                    var list_access_tbody = '';
                    list_access_tbody+='<tr><td>'+user.username+'</td>'
                    accesslist.servers.forEach(function (server) {
                        if(userAccessArray[user.id] && userAccessArray[user.id].includes(server.id)) {
                            list_access_tbody+= '<th><label><input type="checkbox" class="filled-in access-checkbox" checked="checked" data-userid="'+user.id+'" data-serverid="'+server.id+'" /><span></span></label></th>'
                        } else {
                            list_access_tbody+= '<th><label><input type="checkbox" class="filled-in access-checkbox" data-userid="'+user.id+'" data-serverid="'+server.id+'" /><span></span></label></th>'
                        }
                    });
                    list_access_tbody+='</tr>'
                    $("#list_access_tbody").append(list_access_tbody);
                });
            },
            error: function(error) {
                M.toast({html: error.responseText, classes: 'red'});
            }
        })
    })

    $(document.body).on('click','.delete_user_button', function(e) {
        e.preventDefault();
        $.ajax({
            url: '/deleteUser',
            data: {
                'user_id': $(this).data('userid')
            },
            success: function(data) {
                M.toast({html: data, classes: 'green'});
                $("#list_users").click();
            },
            error: function(error) {
                M.toast({html: error.responseText, classes: 'red'});
            }
        })
    })    

    $(document.body).on('click','.delete_server_button', function(e) {
        e.preventDefault();
        $.ajax({
            url: '/deleteServer',
            data: {
                'server_id': $(this).data('serverid')
            },
            success: function(data) {
                M.toast({html: data, classes: 'green'});
                $("#list_servers").click();
            },
            error: function(error) {
                M.toast({html: error.responseText, classes: 'red'});
            }
        })
    })    
    
    $(document.body).on('click','.access-checkbox', function(e) {
        console.log("toggled")
        // e.preventDefault();
        $.ajax({
            url: '/toggleAccess',
            data: {
                'user_id': $(this).data('userid'),
                'server_id': $(this).data('serverid'),
                'access': $(this).prop("checked")
            },
            success: function(data) {
                M.toast({html: data, classes: 'green'});
                // var instance = M.Modal.getInstance($("#add_server_modal"));
                // instance.close();
                // reset();
            },
            error: function(error) {
                M.toast({html: error.responseText, classes: 'red'});
            }
        })
    })

});