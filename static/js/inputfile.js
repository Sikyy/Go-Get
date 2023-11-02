//inputfile.js

const selectFileButton = document.getElementById("select-file-button");
const fileInput = document.getElementById("file-input");
const selectedFileInfo = document.getElementById("selected-file-info");

selectFileButton.addEventListener("click", function() {
    fileInput.click();
});

fileInput.addEventListener("change", function() {
    if (fileInput.files.length > 0) {
        const selectedFile = fileInput.files[0];
        selectedFileInfo.textContent = "你选择的文件名是：" + selectedFile.name
    } else {
        selectedFileInfo.textContent = "没有选择文件";
    }
});
