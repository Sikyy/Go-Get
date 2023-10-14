//task-list.js
const ws = new WebSocket("ws://localhost:9000/ws");

// 处理 WebSocket 连接的事件
ws.onopen = function() {
    console.log("WebSocket 连接已建立");
};

switch (true) {
    case message.startsWith("文件原始名称:"):
        // 处理文件原始名称消息
        const fileName = message.substring("文件原始名称:".length);
        // 添加下载任务到下载任务列表
        addDownloadTask(fileName, ""); // 传入文件名和文件大小
        break;
    case message.startsWith("文件大小:"):
        // 处理文件大小消息
        const fileSize = message.substring("文件大小:".length);
        // 更新下载任务列表中的文件大小
        updateDownloadTaskSize(fileSize);
        break;
    default:
        // 处理其他消息类型
        break;
}

function addDownloadTask(fileName, fileSize) {
const magnetList = document.getElementById("magnet-list");

// 创建新的 <li> 元素
const newItem = document.createElement("li");
const taskName = document.createElement("span");
const taskSize = document.createElement("span");
const progress = document.createElement("span");

taskName.textContent = "任务名称: " + fileName;
taskSize.textContent = "文件大小: " + fileSize;
progress.textContent = "下载进度: 0%";

newItem.appendChild(taskName);
newItem.appendChild(taskSize);
newItem.appendChild(progress);

magnetList.appendChild(newItem);
}

function updateDownloadTaskSize(fileSize) {
// 在下载任务列表中更新文件大小
const fileSizeElement = document.getElementById("file-size-element");
fileSizeElement.textContent = "文件大小: " + fileSize;
}