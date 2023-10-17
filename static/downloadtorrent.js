//downloadtorrent.js 
    
document.getElementById("add-torrent-btn").addEventListener("click", function () {
    // 使用 fetch 发起 GET 请求到后端
    fetch("/torrent", {
        method: "GET"
    })
    .then(response => {
        if (response.ok) {
            // 下载成功
            alert("下载成功！");
        } else {
            // 下载失败
            alert("下载失败！");
        }
    })
    .catch(error => {
        // 捕获其他错误
        console.error("Error:", error);
    });
});