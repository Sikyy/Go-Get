var folderSelector = document.getElementById('folder-selector');
folderSelector.addEventListener('change', handleFolderSelect);

function handleFolderSelect(event) {
    var files = event.target.files;
    if (files.length > 0) {
    var folderPath = files[0].path;
    console.log(folderPath);
    }
}

