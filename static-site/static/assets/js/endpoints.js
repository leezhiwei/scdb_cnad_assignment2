let host = window.location.host // get hostname
let endpoints = {
	"login": `https://${host}/api/v1/login`,
	"video-conf": `https://100.81.233.97:8080/api/v1/videoconf`,
	"medqna": `https://${host}/api/v1/medqna`,
	"gait_analysis": `http://192.168.2.146:8501`,
	"assessment": `https://${host}/api/v1/assessment`,
	"healthguide": `https://${host}/api/v1/healthguide`
}
