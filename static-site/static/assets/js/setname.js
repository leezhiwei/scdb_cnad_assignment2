function getFormData($form){
    // use jquery to serialise form into JSON
    var unindexed_array = $form.serializeArray();
    var indexed_array = {};

    $.map(unindexed_array, function(n, i){
        indexed_array[n['name']] = n['value'];
    });

    return indexed_array;
}
$('#submitbutton').click(function () {
    let username = getFormData($('#nameInput')).peername
    if (username == ""){
        alert("Username cannot be blank.")
        return
    }
    document.cookie = "username=" + username
    window.location.href = "calldoctor.html"
});