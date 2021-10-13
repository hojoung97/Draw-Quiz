/******************* Room Arguments *******************/
let userName
let roomID
new URLSearchParams(window.location.search).forEach((value, name) => {
    if (name == "userName") {
        userName = value
    } else if  (name == "roomID") {
        roomID = value
    }
});

/******************* Websocket code *******************/
let socket = new WebSocket(`ws://localhost:8050/${roomID}/${userName}`);

socket.onopen = () => {
    console.log("Connected to Websocket Server");
};

socket.onmessage = (msg) => {
    let data = JSON.parse(msg.data);

    if (data.type == 1) {
        console.log(data.body)

        s = data.body.split(";");

        if (s[0] == "choose") {
            roomLog.innerHTML = "Ready to begin! Choose an object below.";
            pickItems();
            socket.send("wait");
        } else if (s[0] == "wait") {
            roomLog.innerHTML = "Friend is choosing something to draw";
        } else if (s[0] == "drawing") {
            roomLog.innerHTML = "Friend is drawing the object";
        } else if (s[0] == "done") {
            roomLog.innerHTML = "Can you guess the object?";
            for (let i=0; i < 4; i++) {
                document.getElementById(`option${i+1}`).innerHTML = s[i+1];
            }
            guessItems();
        } else if (s[0] == "correct0") {
            window.alert("You got it right!");
            roomLog.innerHTML = "Ready to begin! Choose an object below.";
            pickItems();
            socket.send("wait");
        } else if (s[0] == "correct1") {
            window.alert("Your friend got it right!");
        } else if (s[0] == "wrong0") {
            window.alert("You got it wrong :(");
            roomLog.innerHTML = "Ready to begin! Choose an object below.";
            pickItems();
            socket.send("wait");
        } else if (s[0] == "wrong1") {
            window.alert("Your friend got in wrong :(");
        }

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
    window.alert("Room is full! Redirecting to lobby...");
    window.location.replace("/");
};

/******************* UI Related *******************/
const welcomeMess = document.getElementById("welcomeMess");
welcomeMess.append(`Room #${roomID}: Welcome ${userName}!`)

const drawBox = document.querySelector(".drawBox");
const roomLogBox = document.querySelector(".roomLogBox");
const selectionBox = document.querySelector(".selectionBox");

const roomLog = document.getElementById("roomLog");
const selectionButton = document.getElementById("selectionButton");
const ansButton = document.getElementById("ansButton");

function pickItems() {
    options = [];
    for (let i=0; i < 4; i++) {
        options.push(document.getElementById(`option${i+1}`));
    }

    items = shuffle(items);
    for (let i=0; i < options.length; i++) {
        options[i].innerHTML = items[i];
        options[i].setAttribute("onclick", "optionSelected(this, options, selectionButton)");
    }

    let opts = `option;${items[0]};${items[1]};${items[2]};${items[3]}`;
    socket.send(opts);
}

function guessItems() {
    options = [];
    for (let i=0; i < 4; i++) {
        options.push(document.getElementById(`option${i+1}`));
    }

    for (let i=0; i < 4; i++) {
        options[i].setAttribute("onclick", "optionSelected(this, options, ansButton)");
    }
}

let selectedItem;
function optionSelected(option, options, button) {
    for (let i=0; i < options.length; i++) {
        options[i].classList.remove("chosen");
    }

    selectedItem = option.textContent;
    console.log(selectedItem);
    option.classList.add("chosen");
    button.style.visibility = "visible";
}

selectionButton.addEventListener("click", e => {
    socket.send(`drawing;${selectedItem}`);
    roomLog.innerHTML = `You chose ${selectedItem}. Start drawing!`;
    selectionButton.style.visibility = "hidden";
});

ansButton.addEventListener("click", e => {
    socket.send(`answer;${selectedItem}`);
    ansButton.style.visibility = "hidden";
});

/******************* Canvas Code *******************/
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
    socket.send("done");
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
