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

      // Create a table for changes
      const table = document.createElement("table");
      table.classList.add("changes-table"); // Add a class for styling

      // Create table header
      const thead = document.createElement("thead");
      const headerRow = document.createElement("tr");

      const thCheckbox = document.createElement("th");
      thCheckbox.textContent = "Select";
      thCheckbox.style.width = "50px"; // Give some fixed width
      thCheckbox.style.textAlign = "center";
      headerRow.appendChild(thCheckbox);

      const thType = document.createElement("th");
      thType.textContent = "Type";
      thType.style.width = "120px"; // Fixed width for type
      headerRow.appendChild(thType);

      const thFile = document.createElement("th");
      thFile.textContent = "File / Directory";
      headerRow.appendChild(thFile); // This column will take remaining space

      const thTarget = document.createElement("th");
      thTarget.textContent = "Target Name";
      headerRow.appendChild(thTarget);

      thead.appendChild(headerRow);
      table.appendChild(thead);

      // Create table body
      const tbody = document.createElement("tbody");

      changes.forEach((change, index) => {
        const tr = document.createElement("tr");
        tr.classList.add("change-item"); // Add a class for row styling

        // Cell for checkbox
        const tdCheckbox = document.createElement("td");
        tdCheckbox.style.textAlign = "center";
        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.checked = change.selected;
        checkbox.addEventListener("change", () => {
          change.selected = checkbox.checked;
        });
        tdCheckbox.appendChild(checkbox);
        tr.appendChild(tdCheckbox);

        // Cell for change type
        const tdType = document.createElement("td");
        let typeText = "";
        if (change.type === "delete_file") {
          typeText = "[DELETE]";
        } else if (change.type === "rename_file") {
          typeText = "[RENAME]";
        } else if (change.type === "remove_dir") {
          typeText = "[REMOVE DIR]";
        } else {
            typeText = "[UNKNOWN]";
        }
        tdType.textContent = typeText;
        tr.appendChild(tdType);

        // Cell for file path (target)
        const tdFile = document.createElement("td");
        const filePathSpan = document.createElement("span"); // Create a span to hold the text
        filePathSpan.classList.add("file-path-text"); // Add a class for styling
        filePathSpan.textContent = change.target || "";
        tdFile.appendChild(filePathSpan); // Append the span to the td
        tr.appendChild(tdFile);

        // Cell for new name (target name for renames)
        const tdNewName = document.createElement("td");
        const newNameSpan = document.createElement("span");
        newNameSpan.classList.add("file-path-text"); // Apply the class for dynamic fading
        newNameSpan.textContent = change.newName || ""; // Set the text content
        tdNewName.appendChild(newNameSpan); 
        tr.appendChild(tdNewName);

        tbody.appendChild(tr);
      });

      table.appendChild(tbody);
      changesList.appendChild(table);
      setTimeout(applyFadeToOverflowingPaths, 0);
    } catch (error) {
      console.error("Error fetching changes:", error);
      changesList.innerHTML =
        "<p class='result'>Error loading changes: " + error + "</p>";
      window.runtime.LogError("Error fetching changes: " + error);
    }
  });
};

/**
 * Dynamically applies or removes the fade effect and ellipsis
 * based on whether the file path text overflows its container.
 * This function should be called after the table is rendered
 * and whenever the window is resized.
 */
function applyFadeToOverflowingPaths() {
    const allFilePathSpans = document.querySelectorAll(".file-path-text");
    allFilePathSpans.forEach(span => {
        // `scrollWidth` is the intrinsic width of the content.
        // `clientWidth` is the visible width of the element.
        // These are accurate because `.file-path-text` has `white-space: nowrap;` and `overflow: hidden;`.

        if (span.scrollWidth > span.clientWidth) {
            // Text overflows: apply the fade and ensure ellipsis is *not* clipped
            span.classList.add("fade-overflow");
            span.style.textOverflow = 'ellipsis'; // Explicitly set, as we'll clip otherwise
        } else {
            // Text fits: remove the fade and ensure ellipsis is clipped
            span.classList.remove("fade-overflow");
            span.style.textOverflow = 'clip'; // Effectively hides the ellipsis if text fits
        }
    });
}

window.addEventListener('resize', applyFadeToOverflowingPaths);