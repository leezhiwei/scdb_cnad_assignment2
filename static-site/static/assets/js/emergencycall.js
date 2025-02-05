    $("#AddNumber").on("click", function () {
        let contactName = $("#name").val().trim();
        let contactNumber = $("#phone").val().trim();
        let seniorId = getCookie("senior_id");

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


    function getCookie(name) {
        let match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
        return match ? match[2] : null;
    }
    
    $("#ListEmergencyNumbers").on("click", function () {
        let seniorId = getCookie("senior_id");
    
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
                let contactsList = JSON.parse(response) || []
                let contactsHtml;

                if (contactsList.length === 0) {
                    contactsHtml = `<div class="alert alert-warning text-center">No emergency contacts found.</div>`;
                } else {
                    contactsHtml = `
                        <table class="table table-bordered table-striped">
                            <thead class="table-dark">
                                <tr>
                                    <th>Name</th>
                                    <th>Phone Number</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${contactsList.map(contact => `
                                    <tr>
                                        <td>${contact.contactname}</td>
                                        <td>${contact.contactnumber}</td>
                                        <td>
                                            <button> <a href="tel:+65${contact.contactnumber}">call this number</a></button>
                                            <button class="btn btn-danger btn-sm delete-contact" data-id="${contact.id}">Delete</button>
                                        </td>
                                    </tr>
                                `).join("")}
                            </tbody>
                        </table>
                    `;
                }
    
                $("#emergencyContacts").html(contactsHtml);
            },
            error: function (xhr) {
                alert("Error fetching emergency contacts: " + xhr.responseText);
            }
        });
    });
    
    
    $(document).on("click", ".delete-contact", function () {
        let contactId = $(this).data("id"); // Get ID from the button's data-id attribute
    
        if (!contactId) {
            alert("Invalid contact ID.");
            return;
        }
    
        $.ajax({
            url: "http://localhost:8080/api/v1/login/deleteemergencycontact",
            type: "POST",
            contentType: "application/json",
            data: JSON.stringify({ contact_id: parseInt(contactId) }),
            success: function (response) {
                alert(response.message || "Emergency contact deleted successfully!");
                location.reload(); // Refresh the list after deletion
            },
            error: function (xhr) {
                alert("Error deleting emergency contact: " + xhr.responseText);
            }
        });
    });