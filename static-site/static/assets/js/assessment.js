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
        // Obtain seniorID
        let seniorId = document.cookie = "senior_id";
        // Prevent page reload
        event.preventDefault(); 

        let formData = getFormData($("#healthasst"))
        formData.senior_id = seniorId;

        $.ajax({
            url: "http://localhost:8081/api/v1/assessment/submit",
            type: "POST",
            data: JSON.stringify(formData),
            contentType: "application/json",
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