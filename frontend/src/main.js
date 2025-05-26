import './style.css';
import './app.css';

import logo from './assets/images/logo.png';
import { SimulateScan } from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
  <img id="logo" class="logo">
  <div class="result" id="result">Enter a directory path to simulate scan:</div>
  <div class="input-box" id="input">
    <input class="input" id="dirInput" type="text" autocomplete="off" placeholder="Enter directory path" />
    <button class="btn" id="scanBtn">Simulate Scan</button>
  </div>
`;
document.getElementById('logo').src = logo;

let inputElement = document.getElementById("dirInput");
inputElement.focus();
let resultElement = document.getElementById("result");

document.getElementById('scanBtn').onclick = async function () {
  let dir = inputElement.value.trim();

  if (dir === "") {
    alert("Please enter a directory path.");
    return;
  }

  resultElement.innerText = "Scanning...";

  try {
    // Call the Go backend function you'll expose (to be implemented)
    let scanResult = await SimulateScan(dir);
    resultElement.innerText = scanResult;
  } catch (err) {
    console.error(err);
    resultElement.innerText = "Error: " + err.message;
  }
};
