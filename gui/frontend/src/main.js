//import './styles/base.css';
import './styles/theme.css';
import './styles/layout.css';
import './styles/utilities.css';
import './styles/components/buttons.css';
import './styles/components/input.css';
import './styles/components/table.css';

// === Constants ===

const DEFAULT_CHANGES_MSG = "<p class='result'>Select a folder to view changes.</p>";
const NO_CHANGES_MSG = "<p class='result'>No changes found.</p>";
const NO_FOLDER_SELECTED = "No folder selected";

// === Event Setup ===

window.onload = () => {
  const selectFolderButton = document.getElementById("select-folder");
  const selectedPathSpan = document.getElementById("selected-path");
  const changesList = document.getElementById("changes-list");
  const logOutput = document.getElementById("log-output");

  changesList.innerHTML = DEFAULT_CHANGES_MSG;

  selectFolderButton?.addEventListener("click", async () => {
    try {
      logOutput && (logOutput.textContent = "Opening folder picker...");

      const defaultPath = "D:\\Coding\\HouseKeeper\\testdata";
      const folderPath = await window.go.main.App.OpenDirectoryDialog("Select Folder to Scan", defaultPath);

      if (!folderPath) {
        selectedPathSpan.textContent = NO_FOLDER_SELECTED;
        changesList.innerHTML = DEFAULT_CHANGES_MSG;
        return;
      }

      selectedPathSpan.textContent = folderPath;

      const changes = await window.go.main.App.GetChanges(
        folderPath,
        "../userconfigs/extensions_to_delete.json",
        "../userconfigs/extension_replacements.json"
      );

      renderChanges(changesList, changes);
    } catch (error) {
      console.error("Error fetching changes:", error);
      changesList.innerHTML = `<p class='result'>Error loading changes: ${error}</p>`;
      window.runtime?.LogError?.("Error fetching changes: " + error);
    }
  });
};

// === Helper Functions ===

const renderChanges = (container, changes) => {
  container.innerHTML = "";

  if (!changes || changes.length === 0) {
    container.innerHTML = NO_CHANGES_MSG;
    return;
  }

  const table = document.createElement("table");
  table.className = "changes-table";

  table.appendChild(createTableHeader());
  table.appendChild(createTableBody(changes));

  container.appendChild(table);
  setTimeout(applyFadeToOverflowingPaths, 0);
};

const createTableHeader = () => {
  const thead = document.createElement("thead");
  const row = document.createElement("tr");

  const headers = [
    { text: "Select", width: "50px", align: "center" },
    { text: "Type", width: "50px" },
    { text: "File / Directory" },
    { text: "Target Name" },
  ];

  headers.forEach(({ text, width, align }) => {
    const th = document.createElement("th");
    th.textContent = text;
    if (width) th.style.width = width;
    if (align) th.style.textAlign = align;
    row.appendChild(th);
  });

  thead.appendChild(row);
  return thead;
};

const createTableBody = (changes) => {
  const tbody = document.createElement("tbody");

  const typeIcons = {
    delete_file: "/icons/delete-file.svg",
    rename_file: "/icons/rename.svg",
    remove_dir: "/icons/delete-folder.svg",
  };

  changes.forEach((change) => {
    const row = document.createElement("tr");
    row.className = "change-item";

    // Checkbox
    const checkboxCell = document.createElement("td");
    checkboxCell.style.textAlign = "center";
    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.checked = change.selected;
    checkbox.addEventListener("change", () => (change.selected = checkbox.checked));
    checkboxCell.appendChild(checkbox);

/*     // Type
    const typeCell = document.createElement("td");
    typeCell.textContent = typeLabels[change.type] || "[UNKNOWN]"; */

    const typeCell = document.createElement("td");

    const iconSrc = typeIcons[change.type];
    if (iconSrc) {
      const img = document.createElement("img");
      img.src = iconSrc;
      img.alt = change.type;
      img.style.width = "24px";  // Adjust size as needed
      img.style.height = "24px";
      typeCell.appendChild(img);
    } else {
      typeCell.textContent = "[UNKNOWN]";
    }

    // File path
    const fileCell = document.createElement("td");
    fileCell.appendChild(createOverflowSpan(change.target));

    // New name
    const targetCell = document.createElement("td");
    targetCell.appendChild(createOverflowSpan(change.newName));

    [checkboxCell, typeCell, fileCell, targetCell].forEach((cell) => row.appendChild(cell));
    tbody.appendChild(row);
  });

  return tbody;
};

/*
const createOverflowSpan = (text = "") => {
  const span = document.createElement("span");
  span.className = "file-path-text";
  span.textContent = text;
  return span;
};
*/
const createOverflowSpan = (text = "") => {
  const outerSpan = document.createElement("span");
  outerSpan.className = "file-path-text";

  const innerSpan = document.createElement("span");
  innerSpan.className = "file-path-inner";
  innerSpan.textContent = text;

  outerSpan.appendChild(innerSpan);
  return outerSpan;
};

/*
const applyFadeToOverflowingPaths = () => {
  document.querySelectorAll(".file-path-text").forEach((span) => {
    const isOverflowing = span.scrollWidth > span.clientWidth;
    span.classList.toggle("fade-overflow", isOverflowing);
    span.classList.toggle("scrolling", isOverflowing);
    span.style.textOverflow = isOverflowing ? "ellipsis" : "clip";

    if (isOverflowing) {
      const overflow = innerSpan.scrollWidth - outerSpan.clientWidth;
      innerSpan.style.setProperty("--scroll-distance", `-${overflow}px`);
    } else {
      innerSpan.style.removeProperty("--scroll-distance");
    }
  });
};
*/

const applyFadeToOverflowingPaths = () => {
  document.querySelectorAll(".file-path-text").forEach((outerSpan) => {
    const innerSpan = outerSpan.querySelector(".file-path-inner");
    if (!innerSpan) return;

    const isOverflowing = innerSpan.scrollWidth > outerSpan.clientWidth;
    outerSpan.classList.toggle("fade-overflow", isOverflowing);
    outerSpan.classList.toggle("scrolling", isOverflowing);
    outerSpan.style.textOverflow = isOverflowing ? "ellipsis" : "clip";

    if (isOverflowing) {
      const overflow = innerSpan.scrollWidth - outerSpan.clientWidth;
      innerSpan.style.setProperty("--scroll-distance", `${overflow}px`);
    } else {
      innerSpan.style.removeProperty("--scroll-distance");
    }
  });
};


window.addEventListener("resize", applyFadeToOverflowingPaths);