import Alpine from 'alpinejs'
import Clipboard from './clipboard';

window.Alpine = Alpine

Alpine.plugin(Clipboard)
Alpine.start()

document.addEventListener("DOMContentLoaded", function () {
  const themeToggle = document.getElementById("themeToggle");
  const htmlElement = document.documentElement;

  // Check for saved theme or system preference
  if (
    localStorage.theme === "dark" ||
    (!("theme" in localStorage) &&
      window.matchMedia("(prefers-color-scheme: dark)").matches)
  ) {
    htmlElement.classList.add("dark");
  } else {
    htmlElement.classList.remove("dark");
  }

  themeToggle.addEventListener("click", function () {
    console.log("switching color mode")
    if (htmlElement.classList.contains("dark")) {
      htmlElement.classList.remove("dark");
      localStorage.theme = "light";
    } else {
      htmlElement.classList.add("dark");
      localStorage.theme = "dark";
    }
  });

  // Prevent form resubmission on page refresh
  if (window.history.replaceState) {
    window.history.replaceState(null, null, window.location.href);
  }
});

/**
 *
 * @param {string} passphrase
 * @param {string} salt
 */
async function encodePassphrase(passphrase, salt) {
  const encodedPassphrase = new TextEncoder().encode(passphrase);

  const initialKey = await crypto.subtle.importKey(
    "raw",
    encodedPassphrase,
    { name: "PBKDF2" },
    false,
    ["deriveKey"]
  );

  return await crypto.subtle.deriveKey(
    { name: "PBKDF2", salt, iterations: 100000, hash: "SHA-256" },
    initialKey,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt", "decrypt"]
  );
}

async function encrypt(content, passphrase) {
  const salt = crypto.getRandomValues(new Uint8Array(16));

  const key = await encodePassphrase(passphrase, salt);

  const iv = crypto.getRandomValues(new Uint8Array(12));

  const contentBytes = new TextEncoder().encode(content);

  const cipher = new Uint8Array(
    await crypto.subtle.encrypt({ name: "AES-GCM", iv }, key, contentBytes)
  );

  return {
    salt: bytesToBase64(salt),
    iv: bytesToBase64(iv),
    cipher: bytesToBase64(cipher),
  };
}

async function decrypt(encryptedData, password) {
  const salt = base64ToBytes(encryptedData.salt);

  const key = await encodePassphrase(password, salt);

  const iv = base64ToBytes(encryptedData.iv);

  const cipher = base64ToBytes(encryptedData.cipher);

  const contentBytes = new Uint8Array(
    await crypto.subtle.decrypt({ name: "AES-GCM", iv }, key, cipher)
  );

  return new TextDecoder().decode(contentBytes);
}

function bytesToBase64(arr) {
  return btoa(Array.from(arr, (b) => String.fromCharCode(b)).join(""));
}

function base64ToBytes(base64) {
  return Uint8Array.from(atob(base64), (c) => c.charCodeAt(0));
}

window.encrypt = encrypt
window.decrypt = decrypt
