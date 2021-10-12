/*
Websocket code
*/
const welcomeMess = document.getElementById("welcomeMess");
let userName
let roomID
new URLSearchParams(window.location.search).forEach((value, name) => {
    if (name == "userName") {
        userName = value
    } else if  (name == "roomID") {
        roomID = value
    }
});
welcomeMess.append(`Room #${roomID}: Welcome ${userName}!`)

let socket = new WebSocket(`ws://localhost:8080/room/${roomID}`);

socket.onopen = () => {
    console.log("Connected to Server");
};

socket.onmessage = (msg) => {
    let data = JSON.parse(msg.data);
    if (data.type == 1) {
        console.log(data.body)
    } else if (data.type == 2) {
        let image = new Image();
        image.onload = function() {
            context.drawImage(image, 0, 0);
        };
        image.src = data.body
    }
}

socket.onclose = (event) => {
    console.log("Socket Closed Connection: ", event);
};

socket.onerror = (error) => {
    console.log("Socket Error: ", error);
};

/*
Canvas Code
*/
let isDrawing = false;
let x = 0;
let y = 0;

const canvas = document.getElementById("myCanvas");
const context = canvas.getContext("2d");
const brushPalette = document.getElementById("brushPalette");
const brushSlider = document.getElementById("brushSlider");
const clearButton = document.getElementById("clearButton");
const submitButton = document.getElementById("submitButton");

context.fillStyle = "white";
context.fillRect(0, 0, canvas.width, canvas.height);

canvas.addEventListener("mousedown", e => {
    x = e.offsetX;
    y = e.offsetY;
    isDrawing = true;
});

canvas.addEventListener("mousemove", e => {
    if (isDrawing === true) {
        drawLine(x, y, e.offsetX, e.offsetY);
        x = e.offsetX;
        y = e.offsetY;
    }
});

window.addEventListener("mouseup", e => {
    if (isDrawing === true) {
        drawLine(x, y, e.offsetX, e.offsetY);
        x = 0;
        y = 0;
        isDrawing = false;
        broadcastDrawing();
    }
});

clearButton.addEventListener("click", e => {
    context.fillStyle = "white";
    context.fillRect(0, 0, canvas.width, canvas.height);
    broadcastDrawing();
});

submitButton.addEventListener("click", e => {
    let drawing = canvas.toDataURL("image/jpeg");
    let enc = new TextEncoder()
    socket.send(enc.encode(drawing));
});

function broadcastDrawing() {
    let drawing = canvas.toDataURL("image/jpeg");
    let enc = new TextEncoder();
    socket.send(enc.encode(drawing));
}

function drawLine(x1, y1, x2, y2) {
    context.beginPath();
    context.strokeStyle = brushPalette.value;
    context.lineWidth = brushSlider.value;
    context.lineCap = "round";
    context.moveTo(x1, y1);
    context.lineTo(x2, y2);
    context.stroke();
    context.closePath();
}
