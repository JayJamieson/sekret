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
