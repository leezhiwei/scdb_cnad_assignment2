// //init for cookies  as logout user if there is not cookies found
// function init(){
// 	if(typeof Cookies.get('senior_id') !== 'undefined'){
// 		$('#loginbutton').text("Logout");
// 		$('#loginbutton').attr("href","./");
// 		$('#loginbutton').click(function(){
// 			Cookies.remove("user_id", {path: "/", sameSite: "lax"});
// 			var delay = 100; 
// 			setTimeout(function(){ location.reload() }, delay);
// 			return false;
// 	});
// 	}
// 	else
// 	{
// 		$('#viewuser').hide();
// 	}
// }
// Send SMS when the "Send SMS" button is clicked
$("#send-sms").on("click", function () {
    let phone = $("#phone").val().trim();
    if (!phone) {
        alert("Please enter your phone number.");
        return;
    }

    $.ajax({
        url: "http://localhost:8080/api/v1/login/sendsms",
        type: "POST",
        contentType: "application/json",
        data: JSON.stringify({ phone: phone }),
        success: function (response) {
            alert(response.message || "SMS sent successfully!");
        },
        error: function (xhr) {
            alert("Error sending SMS: " + xhr.statusText);
        }
    });
});

// Log in when the "Log In" button is clicked
$("#login").on("click", function () {
    let phone = $("#phone").val().trim();
    let smsCode = $("#sms-code").val().trim();

    if (!phone || !smsCode) {
        alert("Please enter both phone number and SMS code.");
        return;
    }

    $.ajax({
        url: "http://localhost:8080/api/v1/login/login",
        type: "POST",
        contentType: "application/json",
        xhrFields: {
         withCredentials: true
        },
        data: JSON.stringify({ phone: phone, smscode: smsCode }),
        success: function (resptext) {
            try{
                response = JSON.parse(resptext)
            }
            catch{
                alert("Internal Server Error")
                alert(`Error: ${resptext}`)
                console.log(resptext)
            }
            if (response.senior_id) {
                alert("Login successful!");
                //document.cookie = `senior_id=${response.senior_id}; path=/`;
                window.location.href = "index.html";
            } else {
                alert("Invalid OTP or phone number.");
            }
        },
        error: function (xhr) {
            alert("Login failed: " + xhr.responseText);
        }
    });
});

$(document).ready(function () {
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
});

$("#ListEmergencyNumbers").on("click", function () {
    //let seniorId = document.cookie = "senior_id"
    let seniorId = 3;

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
