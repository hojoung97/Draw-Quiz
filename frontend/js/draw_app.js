let socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = () => {
    console.log("Connected to Server");
    socket.send("Hello from the client!");
};

socket.onclose = (event) => {
    console.log("Socket Closed Connection: ", event);
};

socket.onerror = (error) => {
    console.log("Socket Error: ", error);
};

let isDrawing = false;
let x = 0;
let y = 0;

const myCanvas = document.getElementById("myCanvas");
const context = myCanvas.getContext("2d");
const brushPalette = document.getElementById("brushPalette");
const brushSlider = document.getElementById("brushSlider");
const clearButton = document.getElementById("clearButton");
const submitButton = document.getElementById("submitButton");

myCanvas.addEventListener("mousedown", e => {
    x = e.offsetX;
    y = e.offsetY;
    isDrawing = true;
});

myCanvas.addEventListener("mousemove", e => {
    if (isDrawing === true) {
        drawLine(context, x, y, e.offsetX, e.offsetY);
        x = e.offsetX;
        y = e.offsetY;
    }
});

window.addEventListener("mouseup", e => {
    if (isDrawing === true) {
        drawLine(context, x, y, e.offsetX, e.offsetY);
        x = 0;
        y = 0;
        isDrawing = false;
    }
});

function drawLine(context, x1, y1, x2, y2) {
    context.beginPath();
    context.strokeStyle = brushPalette.value;
    context.lineWidth = brushSlider.value;
    context.lineCap = "round";
    context.moveTo(x1, y1);
    context.lineTo(x2, y2);
    context.stroke();
    context.closePath();
}


clearButton.addEventListener("click", e => {
    context.clearRect(0, 0, myCanvas.width, myCanvas.height);
});


submitButton.addEventListener("click", e => {
    let drawing = context.getImageData(0, 0, myCanvas.width/2, myCanvas.height);
    context.putImageData(drawing, myCanvas.width/2, 0);
});
