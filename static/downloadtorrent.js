//downloadtorrent.js
const selectFileButton = document.getElementById("select-file-button");
    const selectedFileInput = document.getElementById("file-input");
    const selectedFileInfoDiv = document.getElementById("selected-file-info");
    const addTorrentButton = document.getElementById("add-torrent-btn");

    const uploadDir = "/Users/siky/go/src/Go-Get/static/uploads";

    selectFileButton.addEventListener("click", () => {
      selectedFileInput.click();
    });

    selectedFileInput.addEventListener("change", () => {
      const selectedFile = selectedFileInput.files[0];

      // 获取文件名
      const fileName = selectedFile.name;

      // 获取文件大小
      const fileSize = selectedFile.size;

      // 获取文件类型
      const fileType = selectedFile.type;

      // 获取文件路径

      const filePath = selectedFile.path;

      // 显示文件信息
      selectedFileInfoDiv.innerHTML = `
        <p>文件名：${fileName}</p>
        <p>文件大小：${fileSize} 字节</p>
        <p>文件类型：${fileType}</p>
        <p>文件路径：${filePath}</p>
      `;
    });

    addTorrentButton.addEventListener("click", () => {

      // 获取文件名
      const fileName = selectedFileInput.files[0].name;
      
      // 构造请求数据
      const formData = new FormData();
      formData.append("file", selectedFileInput.files[0], fileName);
  
      // 创建 XMLHttpRequest 对象，用于发送异步 HTTP 请求
      const xhr = new XMLHttpRequest();
  
      // 设置 XMLHttpRequest 对象的 url 属性
      xhr.open("POST", "/upload"); // 发起POST请求到/upload路由
  
      // 处理响应数据
      xhr.onload = () => {
          // 检查响应状态码
          if (xhr.status === 200) {
              // 上传成功
              alert("上传成功");
  
              // 这里可以根据需要执行其他操作
          } else {
              // 上传失败
              alert("上传失败");
          }
      };
  
      // 因为是异步，所以后发送POST请求也没事
      xhr.send(formData);
  });