window.addEventListener("load", () => {
  const hamSidebar = document.getElementById("hamburger-sidebar");
  const hamToggle = document.getElementById("hamburger-menu-toggle");
  const hamInput = document.getElementById("hamburger-menu-input");

  window.addEventListener("click", (event) => {
    if (
      !(hamToggle?.contains(event.target) || hamSidebar?.contains(event.target))
    ) {
      hamInput.checked = false;
    }
  });
});
