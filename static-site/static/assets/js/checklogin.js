let seniorId = getCookie("senior_id");
console.log("Retrieved senior_id:", seniorId);

if (!seniorId) {
    alert("You need to log in first.");
    window.location.href ="login.html"
}