

    $("#AddNumber").on("click", function () {
        let contactName = $("#name").val().trim();
        let contactNumber = $("#phone").val().trim();
        let seniorId = document.cookie = "senior_id";

        if (!contactName || !contactNumber) {
            alert("Please enter both name and phone number.");
            return;
        }

        if (!seniorId) {
            alert("You need to log in first.");
            return;
        }

        $.ajax({
            url: "http://localhost:8080/api/v1/login/addemergencycontact",
            type: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                contactname: contactName,
                contactnumber: contactNumber,
                senior_id: parseInt(seniorId)
            }),
            success: function (response) {
                alert(response.message || "Emergency contact added successfully!");
                $("#emergencyForm")[0].reset(); // Clear the form after success
            },
            error: function (xhr) {
                alert("Error adding emergency contact: " + xhr.responseText);
            }
        });
    });


$("#ListEmergencyNumbers").on("click", function () {
    let seniorId = document.cookie = "senior_id"

    if (!seniorId) {
        alert("You need to log in first.");
        return;
    }

    $.ajax({
        url: "http://localhost:8080/api/v1/login/listemergencycontact",
        type: "POST",
        contentType: "application/json",
        data: JSON.stringify({ senior_id: parseInt(seniorId) }),
        success: function (response) {
            let contactsList = response.contacts || [];

            let contactsHtml = contactsList.length 
                ? contactsList.map(contact => `<li>${contact.contactname} - ${contact.contactnumber}</li>`).join("")
                : "<p>No emergency contacts found.</p>";

            $("#emergencyContacts").html(`<ul>${contactsHtml}</ul>`);
        },
        error: function (xhr) {
            alert("Error fetching emergency contacts: " + xhr.responseText);
        }
    });
});