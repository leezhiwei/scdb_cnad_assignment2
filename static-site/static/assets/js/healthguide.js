$(document).ready(function () {
    // Fetch the Senior ID from cookies
    let seniorId = getCookie("senior_id");

    if (!seniorId) {
        console.error("Senior ID not found in cookies.");
        return;
    }

    // Fetch health guide data from the backend
    $.ajax({
        url: endpoints.healthguide + "/suggestions?senior_id=" + seniorId,
        type: "GET",
        contentType: "application/json",
        success: function (response) {
            if (response.length === 0) {
                $("#healthGuideContainer").html("<p>No health guides available.</p>");
                return;
            }

            // Display each health guide entry in the HTML
            let htmlContent = "";
            response.forEach((guide, index) => {
                let youtubeEmbedURL = guide.healthguide_videolink.replace("watch?v=", "embed/"); // Convert to embed URL
                htmlContent += `
                    <div class="health-guide">
                        <h3>Health Guide Video ${index + 1}</h3>
                        <p>${guide.healthguide_description}</p>
                        <div class="video-container">
                            <iframe src="${youtubeEmbedURL}" frameborder="0" allowfullscreen></iframe>
                        </div>
                        <hr>
                    </div>`;
            });

            // Set to div
            $("#healthGuideContainer").html(htmlContent);
        },
        error: function (xhr, status, error) {
            console.error("Error fetching health guide:", error);
            $("#healthGuideContainer").html("<p>Error loading health guides. Please try again later.</p>");
        }
    });
});

// Function to retrieve cookies by name
function getCookie(name) {
    let match = document.cookie.match(new RegExp("(^| )" + name + "=([^;]+)"));
    return match ? match[2] : null;
}
