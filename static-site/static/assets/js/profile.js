$(document).ready(function () {
    function getCookie(name) {
        let match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
        return match ? match[2] : null;
    }

    let seniorId = getCookie("senior_id");
    console.log("Retrieved senior_id:", seniorId);
    
    if (!seniorId) {
        alert("You need to log in first.");
        return;
    }
    

    $.ajax({
        url: "http://localhost:8080/api/v1/login/getsenior",
        type: "POST",
        data: JSON.stringify({ senior_id: parseInt(seniorId) }),
        success: function (response) {
            let response1 = JSON.parse(response) || []
            $("#name").val(response1.name);
            $("#phone").val(response1.phone);
        },
        error: function (xhr) {
            alert("Error fetching user details: " + xhr.responseText);
        }
    });

    $("#updateProfile").on("click", function () {
        let updatedName = $("#name").val().trim();
        let updatedPhone = $("#phone").val().trim();

        if (!updatedName || !updatedPhone) {
            alert("Please enter both name and phone number.");
            return;
        }

        $.ajax({
            url: "http://localhost:8080/api/v1/login/updatesenior",
            type: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                name: updatedName,
                phone: updatedPhone,
                senior_id: parseInt(seniorId)
            }),
            success: function (response) {
                alert(response.message || "User details updated successfully!");
            },
            error: function (xhr) {
                alert("Error updating user details: " + xhr.responseText);
            }
        });
    });
});
