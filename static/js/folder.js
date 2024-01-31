var folderSelector = document.getElementById('folder-selector');
        folderSelector.addEventListener('change', handleFolderSelect);

        function handleFolderSelect(event) {
            var files = event.target.files;
            if (files.length > 0) {
                var relativePath = files[0].webkitRelativePath;
                var folderPath = relativePath.split('/').slice(0, -1).join('/');
                console.log(folderPath);
            }
        }