//details.js
const ws = new WebSocket("ws://localhost:9000/ws");

// 处理 WebSocket 连接的事件
ws.onopen = function() {
    console.log("WebSocket 连接已建立");
};

ws.onmessage = function(event) {
    const message = event.data;
    const magnetDetailsContainer = document.getElementById("magnet-details-container");

    // 创建一个新的 <p> 元素来显示消息
    const messageElement = document.createElement("p");
    messageElement.textContent = message;

    // 将 <p> 元素添加到 <div id="magnet-details-container">
    magnetDetailsContainer.appendChild(messageElement);
};

ws.onclose = function(event) {
    console.log("WebSocket 连接已关闭");
};