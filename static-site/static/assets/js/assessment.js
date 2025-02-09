function getFormData($form){
    // use jquery to serialise form into JSON
    var unindexed_array = $form.serializeArray();
    var indexed_array = {};

    $.map(unindexed_array, function(n, i){
        indexed_array[n['name']] = n['value'];
    });

    return indexed_array;
}

$(document).ready(function () {
    $("form").on("submit", function (event) {
        event.preventDefault(); // Prevent page reload

        let formData = getFormData($("#healthasst"))

        $.ajax({
            url: "http://localhost:8080/api/v1/assessment/submit",
            type: "POST",
            data: JSON.stringify(formData),
            contentType: "application/json",
            contentType: false,  // Prevent jQuery from setting a default content type
            dataType: "json",  // Ensure the response is treated as JSON
            success: function (response) {
                // Update the content of the #response paragraph
                $("#response").text("Assessment Risk Level: " + response.risk_level);
            },
            error: function (xhr) {
                $("#response").text("Error submitting assessment: " + xhr.responseText);
            }
        });
    });
});