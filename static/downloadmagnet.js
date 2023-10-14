//downloadmagnet.js 

document.getElementById("add-magnet-btn").addEventListener("click", function () {
    const magnetLink = document.getElementById("magnet-link").value;
    
    // 发起 HTTP 请求到后端
    fetch("/magnet?magnetURL=" + magnetLink, {
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