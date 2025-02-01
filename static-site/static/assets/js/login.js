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
                document.cookie = `senior_id=${response.senior_id}; path=/`;
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