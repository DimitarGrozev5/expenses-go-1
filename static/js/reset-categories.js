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

/**
 * Get time periods
 * @typedef {Object} Period
 * @property {number} ID
 * @property {string} Caption
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
    ".used-category-card-template"
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
      unusedHiddenInput.type = "hidden";
      usedHiddenInput.type = "hidden";

      /**
       * @type {Period[]}
       */
      const periods = container?.dataset?.periods.split(";").map((p) => {
        const ps = p.split(",");
        return {
          ID: Number(ps[0]),
          Caption: ps[1],
        };
      });
      if (!periods) return;

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
       * Get free funds
       * @type {HTMLSpanElement | null}
       */
      const freeFundsSpan = container.querySelector(".free-funds");
      if (!freeFundsSpan) return;

      /**
       * Define dynamic values
       */

      /**
       * Free funds value
       * @type {{value: () => number, setValue: (v: number) => void, addEventListener: (cb: (v: number) => void) => void}, removeEventListener: (cb: (v: number) => void) => void}
       */
      const vFreeFunds = Value(Number(freeFundsSpan.textContent));

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
        const unusedElementsKeys = new Set(unusedElements.keys());

        // Loop through unused elements
        for (const elem of unused) {
          // If key is not in allElements
          if (!unusedElementsKeys.has(elem.ID)) {
            // Create new element
            const card = getUnusedElement(elem);

            // Add element to map
            unusedElements.set(elem.ID, card);

            // Add card to DOM
            unusedContainer.appendChild(card);
          } else {
            // Remove key from allElements
            unusedElementsKeys.delete(elem.ID);
          }
        }

        // Loop trough elements that are no longer in unused
        for (const key of unusedElementsKeys) {
          // Remove element from DOM
          unusedElements.get(key).remove();

          // Remove element from map
          unusedElements.delete(key);
        }

        // Add data to input
        unusedHiddenInput.value = catsToText(unused);
      });

      // When used list changes
      vUsedCategories.addEventListener((used) => {
        // Get set of category ids
        const usedElementsKeys = new Set(usedElements.keys());

        // Loop through unused elements
        for (const elem of used) {
          // If key is not in allElements
          if (!usedElementsKeys.has(elem.ID)) {
            // Create new element
            const card = getUsedElement(elem);

            // Add element to map
            usedElements.set(elem.ID, card);

            // Add card to DOM
            usedContainer.appendChild(card);
          } else {
            // Remove key from allElements
            usedElementsKeys.delete(elem.ID);
          }
        }

        // Loop trough elements that are no longer in unused
        for (const key of usedElementsKeys) {
          // Remove element from DOM
          usedElements.get(key).remove();

          // Remove element from map
          usedElements.delete(key);
        }

        // Add data to input
        usedHiddenInput.value = catsToText(used);
      });

      vFreeFunds.addEventListener((val) => {
        freeFundsSpan.textContent = val;
      });

      /**
       * Add initial data
       */
      vUnusedCategories.setValue(initUnusedCats);

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
        card
          .querySelectorAll(".name")
          .forEach((e) => (e.textContent = category.Name));
        card.querySelector(".amount").textContent =
          category.CurrentAmount.toFixed(2);
        card.querySelector(".input").textContent =
          category.BudgetInput.toFixed(2);
        card.querySelector(".spending-limit").textContent =
          category.SpendingLimit.toFixed(2);
        card.querySelector(".period").textContent = category.InputPeriodCaption;

        /**
         * Get dialog
         * @type {HTMLDialogElement}
         */
        const dialog = card.querySelector(".dialog");

        // Add click event listener to card
        card.addEventListener("click", () => {
          // Set dialog values
          dialog.querySelector("[name='add_amount']").value =
            category.BudgetInput;
          dialog.querySelector("[name='budget_input']").value =
            category.BudgetInput;
          dialog.querySelector("[name='spending_limit']").value =
            category.SpendingLimit;
          dialog.querySelector("[name='input_interval']").value =
            category.InputInterval;
          dialog.querySelector("[name='input_period']").value =
            category.InputPeriodId;

          // Open dialog
          dialog.showModal();
        });

        // Block clicks from exiting the dialog
        dialog.addEventListener("click", (e) => e.stopPropagation());

        // Add event listener for dialog close
        dialog.querySelector(".dialog-close").addEventListener("click", (e) => {
          // Stop event propagation
          e.stopPropagation();

          // Close dialog
          dialog.close();
        });

        // Add event listener for dialog submit
        dialog.querySelector("form").addEventListener("submit", (e) => {
          // Prevent form submition
          e.preventDefault();

          // Get form values
          const values = new FormData(dialog.querySelector("form"));

          // Create new category object
          const c = { ...category };
          c.InitialAmount = Number(values.get("add_amount"));
          c.BudgetInput = Number(values.get("budget_input"));
          c.SpendingLimit = Number(values.get("spending_limit"));
          c.InputInterval = Number(values.get("input_interval"));
          c.InputPeriodId = Number(values.get("input_period"));

          // Recalculate free funds
          const free = vFreeFunds.value() - c.InitialAmount;

          // If free funds gets bellow zero
          if (free < 0) {
            // Block execution
            alert("Not enough free funds");
            return;
          }

          // Update free funds
          vFreeFunds.setValue(free);

          // Remove from unused
          vUnusedCategories.setValue(
            vUnusedCategories.value().filter((cat) => cat.ID !== c.ID)
          );

          // Add to used
          vUsedCategories.setValue([...vUsedCategories.value(), c]);

          // Close dialog
          dialog.close();
        });

        return card;
      }

      /**
       * Create new Used Category element
       * @param {CategoryOverview} category
       * @returns {HTMLDivElement}
       */
      function getUsedElement(category) {
        /**
         * Clone card
         * @type {HTMLDivElement | null}
         */
        const card = usedCatgoryTemplate.cloneNode(true).firstChild;
        if (!card) return;

        // Set props
        card.querySelector(".name").textContent = category.Name;
        card.querySelector(".amount-current").textContent = (
          category.CurrentAmount + category.InitialAmount
        ).toFixed(2);
        card.querySelector(".amount-add").textContent =
          category.InitialAmount.toFixed(2);
        card.querySelector(".input").textContent =
          category.BudgetInput.toFixed(2);
        card.querySelector(".spending-limit").textContent =
          category.SpendingLimit.toFixed(2);
        card.querySelector(".period").textContent = `${
          category.InputInterval
        } ${periods.find((p) => p.ID == category.InputPeriodId)?.Caption}`;

        // Add click event listener
        card.addEventListener("click", () => {
          // Calculate old free funds value
          const free = vFreeFunds.value() + category.InitialAmount;

          // Update free funds
          vFreeFunds.setValue(free);

          // Remove from used
          vUsedCategories.setValue(
            vUsedCategories.value().filter((c) => c.ID !== category.ID)
          );

          // Add to unused
          vUnusedCategories.setValue([
            ...vUnusedCategories.value(),
            ...initUnusedCats.filter((c) => c.ID === category.ID),
          ]);
        });

        return card;
      }
    }
  );

  /**
   * Helper functions
   */

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

  /**
   * Map categories to text
   *
   * @param {CategoryOverview[]} data
   * @returns {string}
   */
  function catsToText(data) {
    return data
      .map((c) =>
        [
          c.ID,
          c.Name,
          c.BudgetInput,
          c.InputInterval,
          c.InputPeriodId,
          c.InputPeriodCaption,
          c.SpendingLimit,
          c.PeriodStart,
          c.PeriodEnd,
          c.InitialAmount,
          c.CurrentAmount,
        ].join(",")
      )
      .join(";");
  }
});
