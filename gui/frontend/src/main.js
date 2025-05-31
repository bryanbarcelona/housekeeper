window.onload = function () {
  const selectFolderButton = document.getElementById("select-folder");
  const selectedPathSpan = document.getElementById("selected-path");
  const changesList = document.getElementById("changes-list");

  // Initialize with no changes
  changesList.innerHTML = "<p class='result'>Select a folder to view changes.</p>";

  selectFolderButton.addEventListener("click", async () => {
    try {
      // Open folder picker
      const logOutput = document.getElementById("log-output");
      if (logOutput) logOutput.textContent = "Opening folder picker...";

      const folderPath = await window.go.main.App.OpenDirectoryDialog(
      "Select Folder to Scan",
      "D:\\Coding\\HouseKeeper\\testdata"
      );

      if (!folderPath) {
        selectedPathSpan.textContent = "No folder selected";
        changesList.innerHTML = "<p class='result'>Select a folder to view changes.</p>";
        return;
      }

      // Update selected path display
      selectedPathSpan.textContent = folderPath;

      // Fetch changes
      const changes = await window.go.main.App.GetChanges(
      folderPath,
      "../userconfigs/extensions_to_delete.json",
      "../userconfigs/extension_replacements.json"
      );

      // Clear previous changes
      changesList.innerHTML = "";

      if (changes.length === 0) {
        changesList.innerHTML = "<p class='result'>No changes found.</p>";
        return;
      }

      // Create list of changes with checkboxes
      const ul = document.createElement("ul");
      ul.style.listStyle = "none";
      ul.style.padding = "0";

      changes.forEach((change, index) => {
        const li = document.createElement("li");
        li.className = "result";
        li.style.marginBottom = "8px";

        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.checked = change.selected;
        checkbox.style.marginRight = "8px";
        checkbox.addEventListener("change", () => {
          change.selected = checkbox.checked;
        });

        console.log("Change:", JSON.stringify(change));
        document.getElementById("log-output").textContent = "Change: " + JSON.stringify(change);

        let text = "";
        if (change.type === "delete_file") {
          text = `[DELETE] ${change.target}`;
        } else if (change.type === "rename_file") {
          text = `[RENAME] ${change.target} → ${change.newName}`;
        } else if (change.type === "remove_dir") {
          text = `[REMOVE DIR] ${change.target}`;
        }
        if (!text) {
          text = `[UNKNOWN] ${JSON.stringify(change)}`;
        }

        li.appendChild(checkbox);
        li.appendChild(document.createTextNode(text));
        ul.appendChild(li);
      });

      changesList.appendChild(ul);
    } catch (error) {
      console.error("Error fetching changes:", error);
      changesList.innerHTML =
        "<p class='result'>Error loading changes: " + error + "</p>";
      window.runtime.LogError("Error fetching changes: " + error);
    }
  });
};