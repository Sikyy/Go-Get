// output.js
const outputDiv = document.getElementById('output');
const socket = new WebSocket('ws://localhost:9000/ws');

socket.onmessage = function(event) {
    const message = event.data;
    
    // 创建一个包含时间戳和消息的段落元素
    const messageElement = document.createElement('p');
    messageElement.style.color = 'white'; // 设置文本颜色为白色
    const timestamp = new Date().toLocaleTimeString();
    messageElement.innerHTML = `<span class="timestamp">${timestamp}:</span> ${message}`;
    
    // 添加样式，例如不同的背景颜色
    if (message.includes('下载成功')) {
        messageElement.style.backgroundColor = 'lightgreen';
    } else if (message.includes('下载失败')) {
        messageElement.style.backgroundColor = 'pink';
    }
    
    // 添加消息到容器
    outputDiv.appendChild(messageElement);
    
    // 自动滚动到底部
    outputDiv.scrollTop = outputDiv.scrollHeight;
};
