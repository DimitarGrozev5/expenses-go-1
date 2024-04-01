/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../templates/**/*.htm", "./**/*.css", "../static/html/**/*.htm"],
  theme: {
    extend: {
      colors: {
        "primary-50": "#fafafa",
        "primary-100": "#f4f4f5",
        "primary-200": "#e4e4e7",
        "primary-300": "#d4d4d8",
        "primary-400": "#a1a1aa",
        "primary-500": "#71717a",
        "primary-600": "#52525b",
        "primary-700": "#3f3f46",
        "primary-800": "#27272a",
        "primary-900": "#18181b",
        "primary-950": "#09090b",
      },
    },
  },
  plugins: [],
};
