window.addEventListener("load", () => {
  // Get all toggle-dialog-prompt
  const toggleDialogPrompts = document.querySelectorAll(
    ".toggle-dialog"
  );

  toggleDialogPrompts.forEach((toggleDialogPrompt) => {
    // Get dialog
    const dialog = toggleDialogPrompt.nextElementSibling;

    toggleDialogPrompt.addEventListener("click", () => {
      dialog.showModal();
    });

    // Get dialog close buttons
    const closeButtons = dialog.querySelectorAll(".dialog-close");
    closeButtons.forEach((closeButton) => {
      closeButton.addEventListener("click", () => {
        dialog.close();
      });
    });
  });
});
