window.addEventListener("load", () => {
  // Test to see if the browser supports the HTML template element by checking
  // for the presence of the template element's content attribute.
  if (!("content" in document.createElement("template"))) {
    return;
  }

  // Get tag templates
  /** @type {HTMLDivElement | undefined} */
  const tagsInputTemplate =
    document.querySelector(".js-field-template")?.content;

  /** @type {HTMLDivElement | undefined} */
  const tagTemplate = document.querySelector(".tag-template")?.content;

  /** @type {HTMLDivElement | undefined} */
  const selectedTagTemplate = document.querySelector(
    ".selected-tag-template"
  )?.content;

  // Exit if templates are undefined
  if (!(tagsInputTemplate && tagTemplate && selectedTagTemplate)) {
    return;
  }

  // Get all tags inputs and loop trough them
  document.querySelectorAll(".tag-input").forEach((container) => {
    /**
     * Get page elements
     */

    /**
     * Get text input and hide it
     * @type {HTMLInputElement | null}
     */
    const hiddenInput = container.querySelector(".hidden-tag-field");
    if (!hiddenInput) return;

    /**
     * Get tags elements
     * @type {NodeListOf<HTMLDivElement>}
     */
    const tagsElements = container.querySelectorAll(".all-tags > div");

    /**
     * Get template elements
     */

    /**
     * Get tag input template
     * @type {HTMLDivElement | null}
     */
    const tagsInput = tagsInputTemplate.cloneNode(true).firstChild;
    if (!tagsInput) return;

    /**
     * Get input button
     * @type {HTMLButtonElement | null}
     */
    const addTagButton = tagsInput.querySelector(".tag-input-submit");
    if (!addTagButton) return;

    /**
     * Get tag input
     * @type {HTMLInputElement | null}
     */
    const addTagInput = tagsInput.querySelector(".tag-input-field");
    if (!addTagInput) return;

    /**
     * Get selected tags contianer
     * @type {HTMLDivElement | null}
     */
    const selectedTagsContainer = tagsInput.querySelector(".selected-tags");
    if (!selectedTagsContainer) return;

    /**
     * Define dynamic values
     */

    // Add tag input
    const vAddTagInput = value("");

    /**
     * Selected tags init value
     * @type {Set<string>}
     */
    const selectedTagsInitValue = new Map();

    // Selected tags
    const vSelectedTags = value(selectedTagsInitValue);

    /**
     * Define new tags map
     * @type {Map<string, HTMLDivElement>}
     */
    const newTags = new Map();

    /**
     * Add input to page
     */
    container.insertBefore(tagsInput, hiddenInput);

    // Set field to hidden
    hiddenInput.type = "hidden";

    /**
     * Setup event listeners for page elements
     */

    // Define function for adding new tag, taking the value from state
    function addNewTagFromState() {
      // Get new tag value
      const newTag = vAddTagInput.value();

      // Add new tag
      addNewTag(newTag);
    }

    /**
     * Define function for adding new tag
     * @param {string} newTag
     */
    function addNewTag(newTag) {
      // If the tag is too short, exit
      if (newTag.length < 3) return;

      // Get selected tags
      const selectedTags = new Set(vSelectedTags.value());

      // Add new tag
      selectedTags.add(newTag);

      // Set selected tags
      vSelectedTags.setValue(selectedTags);

      // Cleat new tag value
      vAddTagInput.setValue("");
    }

    // Typing in Add Tag input
    addTagInput.addEventListener("input", (e) => {
      // Update dynamic value
      vAddTagInput.setValue(e.target.value);
    });

    // Pressing Enter in Add Tag input
    addTagInput.addEventListener("keydown", (e) => {
      // If button is not Enter, exit
      if (e.key !== "Enter") return;

      // Prevent Default
      e.preventDefault();

      // Add new tag
      addNewTagFromState();
    });

    // Clicking on Add Tag button
    addTagButton.addEventListener("click", (e) => {
      // Add new tag
      addNewTagFromState();
    });

    // Clicking on a Tag from the All Tags list
    for (const tag of tagsElements) {
      tag.addEventListener("click", (e) => {
        addNewTag(tag.textContent);
      });
    }

    /**
     * Set state change effects
     */

    /**
     * Remove tag
     * @param {string} tag
     */
    function removeTag(tag) {
      // Get values
      const v = new Set(vSelectedTags.value());

      // Remove tag from values
      v.delete(tag);

      // Get tag element
      const tagElement = newTags.get(tag);

      // If tag found
      if (!!tagElement) {
        // Remove event listener
        tagElement.removeEventListener("click", removeTagHandler);

        // Remove from DOM
        tagElement.remove();
      }

      // Remove from new tags
      newTags.delete(tag);

      // Update state
      vSelectedTags.setValue(v);
    }

    /**
     * Remove tag handler
     * @param {PointerEvent} e
     */
    function removeTagHandler(e) {
      removeTag(e.target.textContent);
    }

    // Handle updating the Add tag input
    vAddTagInput.addEventListener((value) => {
      addTagInput.value = value;
    });

    // Handle Add Tag button visibility
    vAddTagInput.addEventListener((value) => {
      addTagButton.style.visibility = value.length > 2 ? "visible" : "hidden";
    });

    // Handle tags filter
    vAddTagInput.addEventListener((value) => {
      // Get value length
      const len = value.length;

      // Loop trough tags
      for (const tag of tagsElements) {
        // If word is too short show all tags
        if (len < 3) {
          tag.style.display = "block";
        }

        // If word is long enough
        else {
          // If text is found, show tag
          if (
            tag.innerText
              .toLocaleLowerCase()
              .search(value.toLocaleLowerCase()) >= 0
          ) {
            tag.style.display = "block";
          }

          // If text is not found hide the tag
          else {
            tag.style.display = "none";
          }
        }
      }
    });

    // Handle updating the hidden field
    vSelectedTags.addEventListener((value) => {
      hiddenInput.value = [...value].join(",");
    });

    // Set focus on tag input
    vSelectedTags.addEventListener(() => {
      addTagInput.focus();
    });

    // Handle changes in selected tags
    vSelectedTags.addEventListener((value) => {
      // Get a copy of the New tags
      const nTags = new Map(newTags);

      // Loop trough value
      for (const tag of value) {
        // If the tag is already added
        if (nTags.has(tag)) {
          // Remove it from new tags
          nTags.delete(tag);

          // Go to next item
          continue;
        }

        /**
         * If the tag is new create an element for it
         * @type {HTMLDivElement}
         */
        const tagElement = selectedTagTemplate.cloneNode(true).firstChild;

        // Set inner text
        tagElement.textContent = tag;

        // Set event listener to remove the tag
        tagElement.addEventListener("click", removeTagHandler);

        // Add tag to DOM
        selectedTagsContainer.appendChild(tagElement);

        // Add to new tags
        newTags.set(tag, tagElement);
      }

      // The elements left in nTags have to be removed
      for (const tag of nTags) {
        // Remove tag
        removeTag(tag);
      }
    });

    /**
     * Set initial values
     */

    // Set value of selected tags
    vSelectedTags.setValue(
      new Set(
        hiddenInput.value.length === 0 ? [] : hiddenInput.value.split(/,\s*/)
      )
    );
  });
});

/**
 * @function
 * @template T
 * @param {T} initValue
 * @returns {{value: () => T, setValue: (v: T) => void, addEventListener: (cb: (v: T) => void) => void}, removeEventListener: (cb: (v: T) => void) => void}
 */
function value(initValue) {
  let v = initValue;

  /** @type {Set<(v: T) => void>} */
  const callbacks = new Set();

  const obj = {
    value() {
      return v;
    },
    setValue(value) {
      callbacks.forEach((cb) => cb(value));
      v = value;
    },
    addEventListener(cb) {
      callbacks.add(cb);
    },
    removeEventListener(cb) {
      callbacks.delete(cb);
    },
  };

  return obj;
}
