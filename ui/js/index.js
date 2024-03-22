window.onload = () => {
  const createBtn = document.getElementById("create");
  const secretInput = document.getElementById("secret");

  secretInput.addEventListener("input", (ev) => {
    const secretLink = document.querySelector(".show")

    if (ev.value === "") {
      return;
    }

    if (secretLink == null) {
      return;
    }

    secretLink.className = "hide";
  });

  createBtn.addEventListener("click", (ev) => {
    fetch("/api/secret", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        secret: secretInput.value,
      }),
    })
      .then((res) => res.json())
      .then((secret) => {
        const secretlink = document.getElementById("secretlink");

        secretInput.value = "";
        secretlink.textContent = `${window.location.href}secret/${secret.name}`;
        document.querySelector(".hide").className = "show"
      });
  });
};
