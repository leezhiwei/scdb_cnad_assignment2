const localVideo = document.getElementById('senderVideo');
const remoteVideo = document.getElementById('receiverVideo');

function getCookie(name) {
  var nameEQ = name + "=";
  var ca = document.cookie.split(';');
  for (var i = 0; i < ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
}

let pcSender = null;
let pcReceiver = null;
let localStream = null;
const meetingId = "testMeeting"; // Replace with your meeting ID logic
const userId = getCookie("username"); 
let peerId = ""; 

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function joinmeeting(){
  await fetch(`https://100.81.233.97:8080/join/${meetingId}/${userId}`,{
    method: "POST"
  });
  while (true){
    if (peerId != ""){
      break;
    }
    fetch(`https://100.81.233.97:8080/peer/${meetingId}/${userId}`,{
      method: "GET"
    }).then(function(response) { return response.json(); })
    .then(function(resp) {
      if (resp.peerID == ""){
        $('#waiting').text("Waiting...")
      }
      else{
        $('#waiting').text("Connected!")
        peerId = resp.peerID
        return
      }
    });
    await sleep(2000);
  }
}


async function init() {
  await joinmeeting();
  pcSender = new RTCPeerConnection({
    iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
  });
  
  pcReceiver = new RTCPeerConnection({
    iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
  });

  pcSender.onicecandidate = event => {
    if (event.candidate === null) {
       fetch(`https://100.81.233.97:8080/webrtc/sdp/m/${meetingId}/c/${userId}/p/${peerId}/s/true`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sdp: btoa(JSON.stringify(pcSender.localDescription)) }) // Base64 encode SDP
      }).then(function(response) { return response.json(); })
      .then(function (data) { console.log(JSON.parse(atob(data.Sdp))); pcSender.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(data.Sdp))))})
      
    }
  };

  pcReceiver.onicecandidate = event => {
    if (event.candidate === null) {
      fetch(`https://100.81.233.97:8080/webrtc/sdp/m/${meetingId}/c/${userId}/p/${peerId}/s/false`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sdp: btoa(JSON.stringify(pcReceiver.localDescription)) }) // Base64 encode SDP
      }).then(function(response) { return response.json(); })
      .then(function (data) { console.log(JSON.parse(atob(data.Sdp))); pcReceiver.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(data.Sdp)))) })
    }
  };
}

async function startWebcam() {
  await init();
  try {
    (navigator.mediaDevices.getUserMedia({ video: true, audio: true })).then((stream) => {
      localVideo.srcObject = stream;
      let tracks = stream.getTracks();
      for (var i = 0; i < tracks.length; i++) {
        pcSender.addTrack(stream.getTracks()[i]);
      }
      pcSender.createOffer().then(d => pcSender.setLocalDescription(d))
    })

    pcSender.addEventListener('connectionstatechange', event => {
      if (pcSender.connectionState === 'connected'){
        console.log("connected")
      }
    });

    pcReceiver.addTransceiver('video', {direction: 'recvonly'})

    pcReceiver.createOffer().then(d => pcReceiver.setLocalDescription(d))

    pcReceiver.ontrack = function (event) {
      remoteVideo.srcObject = event.streams[0]
      remoteVideo.autoplay = true
      remoteVideo.controls = true
    }

  } catch (error) {
    console.error('Error accessing webcam:', error);
  }
}


if (userId == "" || userId == null){
  alert("Username not set, please set it.")
  window.location.href = "Startcall.html"
}
startWebcam();
