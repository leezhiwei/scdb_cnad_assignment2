let host = window.location.host // get hostname
let endpoints = {
	"login": `https://${host}:8080/api/v1/login`,
	"video-conf": `https://${host}:8080/api/v1/videoconf`,
	"medqna": `https://${host}:5000/api/v1/medqna`,
	"gait_analysis": `http://${host}`,
	"assessment": `https://${host}:5000/api/v1/assessment`
}