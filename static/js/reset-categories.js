/**
 * @typedef {Object} CategoryOverview
 * @property {number} ID
 * @property {string} Name
 * @property {number} BudgetInput
 * @property {number} InputInterval
 * @property {number} InputPeriodId
 * @property {string} InputPeriodCaption
 * @property {number} SpendingLimit
 * @property {number} PeriodStart
 * @property {number} PeriodEnd
 * @property {number} InitialAmount
 * @property {number} CurrentAmount
 */

window.addEventListener("load", () => {
  // Test to see if the browser supports the HTML template element by checking
  // for the presence of the template element's content attribute.
  if (!("content" in document.createElement("template"))) {
    return;
  }

  // Get templates
  /** @type {HTMLDivElement | undefined} */
  const unusedCatgoryTemplate = document.querySelector(
    ".unused-category-card-template"
  )?.content;

  const usedCatgoryTemplate = document.querySelector(
    ".unused-category-card-template"
  )?.content;

  // Exit if templates are undefined
  if (!(unusedCatgoryTemplate && usedCatgoryTemplate)) {
    return;
  }

  // Get all category-reset forms and loop trough them
  [...document.querySelectorAll(".categories-reset-form")].forEach(
    (container) => {
      /**
       * Get page elements
       */

      /**
       * Get text inputs and hide them
       * @type {HTMLInputElement | null}
       */
      const unusedHiddenInput = container.querySelector(
        ".hidden-unused-categories"
      );
      const usedHiddenInput = container.querySelector(
        ".hidden-used-categories"
      );
      if (!(unusedHiddenInput && usedHiddenInput)) return;
      // unusedHiddenInput.type = "hidden";
      // usedHiddenInput.type = "hidden";

      /**
       * Get unused categories container
       * @type {HTMLDivElement | null}
       */
      const unusedContainer = container.querySelector(
        ".unused-categories-container"
      );
      if (!unusedContainer) return;

      /**
       * Get unused categories container
       * @type {HTMLDivElement | null}
       */
      const usedContainer = container.querySelector(
        ".used-categories-container"
      );
      if (!usedContainer) return;

      /**
       * Define dynamic values
       */

      /**
       * Unused categories values
       * @type {CategoryOverview[]}
       */
      const initUnusedCats = textToCats(unusedHiddenInput.value);

      /**
       * Unused categories list
       * @type {{value: () => CategoryOverview[], setValue: (v: CategoryOverview[]) => void, addEventListener: (cb: (v: CategoryOverview[]) => void) => void}, removeEventListener: (cb: (v: CategoryOverview[]) => void) => void}
       */
      const vUnusedCategories = Value([]);

      /**
       * Used categories values
       * @type {CategoryOverview[]}
       */
      const initUsedCats = textToCats(usedHiddenInput.value);

      /**
       * Used categories list
       * @type {{value: () => CategoryOverview[], setValue: (v: CategoryOverview[]) => void, addEventListener: (cb: (v: CategoryOverview[]) => void) => void}, removeEventListener: (cb: (v: CategoryOverview[]) => void) => void}
       */
      const vUsedCategories = Value([]);

      /**
       * Define rendered elements
       */

      /**
       * Unused elements
       * @type {Map<number, HTMLDivElement>}
       */
      const unusedElements = new Map();

      /**
       * Unused elements
       * @type {Map<number, HTMLDivElement>}
       */
      const usedElements = new Map();

      /**
       * Add effects
       */

      // When unused list changes
      vUnusedCategories.addEventListener((unused) => {
        // Get set of category ids
        const allElements = new Set(unusedElements.keys());

        // Loop through unused keys
        for (const elem of unused) {
          // If key is not in allElements
          if (!allElements.has(elem.ID)) {
            // Create new element
            const card = getUnusedElement(elem);

            // Add element to map
            unusedElements.set(elem.ID, card);

            // Add card to DOM
            unusedContainer.appendChild(card);
          } else {
            // Remove key from allElements
            allElements.delete(elem.ID);
          }
        }

        // Loop trough elem
      });

      /**
       * Add initial data
       */
      vUnusedCategories.setValue(initUnusedCats);
    }
  );

  /**
   * Helper functions
   */

  /**
   * Create new Unused Category element
   * @param {CategoryOverview} category
   * @returns {HTMLDivElement}
   */
  function getUnusedElement(category) {
    /**
     * Clone card
     * @type {HTMLDivElement | null}
     */
    const card = unusedCatgoryTemplate.cloneNode(true).firstChild;
    if (!card) return;

    // Set props
    card.querySelector(".name").textContent = category.Name;
    card.querySelector(".amount").textContent =
      category.CurrentAmount.toFixed(2);
    card.querySelector(".input").textContent = category.BudgetInput.toFixed(2);
    card.querySelector(".spending-limit").textContent =
      category.SpendingLimit.toFixed(2);
    card.querySelector(".period").textContent = category.InputPeriodCaption;

    return card;
  }

  /**
   * Map text to categories
   *
   * @param {string} data
   * @returns {CategoryOverview[]}
   */
  function textToCats(data) {
    return data.split(";").flatMap((cat) => {
      // Exit if empty string
      if (cat.length === 0) return [];

      /** @type {CategoryOverview} */
      const c = {};

      const props = cat.split(",");

      // Check length of props
      if (props.length < 11) {
        console.error(
          "Not enough properties passed for Category Overview in Reset Categories Form"
        );
        return;
      }

      // Add props
      c.ID = Number(props[0]);
      c.Name = props[1];
      c.BudgetInput = Number(props[2]);
      c.InputInterval = Number(props[3]);
      c.InputPeriodId = Number(props[4]);
      c.InputPeriodCaption = props[5];
      c.SpendingLimit = Number(props[6]);
      c.PeriodStart = Number(props[7]);
      c.PeriodEnd = Number(props[8]);
      c.InitialAmount = Number(props[9]);
      c.CurrentAmount = Number(props[10]);

      return c;
    });
  }
});
