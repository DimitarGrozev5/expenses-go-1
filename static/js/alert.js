
/**
 * 
 * @param {string} text 
 * @param {'flash' | 'warn' | 'error'} type 
 * @returns 
 */
async function flashAlert(text, type) {
  // Exit if message is empty
  if (!text) return;

  // Get alert elements
  const alertMsg = document.querySelector(".alert-msg");
  const alertText = document.querySelector(".alert-msg > p");
  const alertClose = document.querySelector(".alert-msg > button");

  // Set type
  switch (type) {
    case "flash":
      alertMsg.classList.remove("warn");
      alertMsg.classList.remove("error");
      break;

    case "warn":
      alertMsg.classList.add("warn");
      alertMsg.classList.remove("error");
      break;

    case "error":
      alertMsg.classList.remove("warn");
      alertMsg.classList.add("error");
      break;

    default:
      alertMsg.classList.remove("warn");
      alertMsg.classList.remove("error");
      break;
  }

  // Set handler functions
  /** @type {number} */
  let timer;
  const openMsg = () => {
    alertMsg.classList.remove("closed");

    clearTimeout(timer);

    timer = setTimeout(() => {
      closeMsg();

      alertClose.removeEventListener("click", closeMsg);
    }, 3000);
  };
  const closeMsg = () => {
    clearTimeout(timer);
    alertMsg.classList.add("closed");
    alertClose.removeEventListener("click", closeMsg);
  };

  // Set button handler
  alertClose.addEventListener("click", closeMsg);

  alertText.innerText = text;
  openMsg();
}
