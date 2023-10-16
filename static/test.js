// 获取输出的div元素
const outputDiv = document.getElementById('output');
const socket = new WebSocket('ws://localhost:9000/ws');

socket.onmessage = function(event) {
    outputDiv.innerHTML += '<p>' + event.data + '</p>';
};