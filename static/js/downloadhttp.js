//downloadmagnet.js 

document.getElementById("add-http-btn").addEventListener("click", function () {
    const httpLink = document.getElementById("http-link").value;
    
    // 发起 HTTP 请求到后端
    fetch("/download?downloadUrl=" + httpLink, {
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