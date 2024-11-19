document.addEventListener("DOMContentLoaded", () => {
    const createFolderForm = document.getElementById("create-folder-form");
    const uploadFileForm = document.getElementById("upload-file-form");
    const folderList = document.getElementById("folder-list");

    // Fetch and display the folder contents
    fetchFolderContents();

    // Handle folder creation
    createFolderForm.addEventListener("submit", async (event) => {
        event.preventDefault();

        const folderName = document.getElementById("folder-name").value;

        const response = await fetch("/folder/create", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                folder_name: folderName,
            }),
        });

        if (response.ok) {
            alert("Folder created successfully!");
            document.getElementById("folder-name").value = ""; // Clear input
            fetchFolderContents(); // Refresh the folder list
        } else {
            const errorData = await response.json();
            alert(`Error creating folder: ${errorData.error}`);
        }
    });

    // Handle file upload
    uploadFileForm.addEventListener("submit", async (event) => {
        event.preventDefault();

        const fileInput = document.getElementById("file-upload");
        const formData = new FormData();
        formData.append("file", fileInput.files[0]);

        const response = await fetch("/file/upload", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            alert("File uploaded successfully!");
            fileInput.value = ""; // Clear input
            fetchFolderContents(); // Refresh the folder list
        } else {
            const errorData = await response.json();
            alert(`Error uploading file: ${errorData.error}`);
        }
    });

    // Fetch folder contents
    async function fetchFolderContents() {
        const response = await fetch("/api/folders");

        if (response.ok) {
            const data = await response.json();
            renderFolderList(data);
        } else {
            alert("Failed to fetch folder contents.");
        }
    }

    // Render the folder and file list dynamically
    function renderFolderList(items) {
        folderList.innerHTML = ""; // Clear existing list

        items.forEach((item) => {
            const listItem = document.createElement("li");
            listItem.textContent = `${item.name} (${item.type})`;

            // Add share button for files
            if (item.type === "file") {
                const shareButton = document.createElement("button");
                shareButton.textContent = "Share";
                shareButton.addEventListener("click", () => handleShare(item.id));
                listItem.appendChild(shareButton);
            }

            folderList.appendChild(listItem);
        });
    }

    // Handle sharing of a file
    async function handleShare(fileId) {
        const accessType = prompt("Enter access type (read/write):", "read");
        const expiration = prompt("Enter expiration date (YYYY-MM-DDTHH:MM:SS) or leave blank:");
        const password = prompt("Enter password (optional):");

        const response = await fetch("/file/share", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                file_id: fileId,
                access_type: accessType,
                expiration: expiration,
                password: password,
            }),
        });

        if (response.ok) {
            const data = await response.json();
            const shareLink = `${window.location.origin}/file/share/${data.share_link}`;
            alert(`Shareable link: ${shareLink}`);
            console.log(`Shareable link: ${shareLink}`);
        } else {
            const error = await response.json();
            alert(`Error sharing file: ${error.error}`);
        }
    }
});
