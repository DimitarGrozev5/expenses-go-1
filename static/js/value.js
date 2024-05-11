/**
 * @function
 * @template T
 * @param {T} initValue
 * @returns {{value: () => T, setValue: (v: T) => void, addEventListener: (cb: (v: T) => void) => void}, removeEventListener: (cb: (v: T) => void) => void}
 */
function Value(initValue) {
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
