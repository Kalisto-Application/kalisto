/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    colors: {
      transparent: "transparent",
      // TODO: Discuss with Designer meaningful names for colors (e.g. primary, secondary, etc.) with theming in mind
      layoutBorder: "#393A4E",
      buttonPrimary: "#2F2EE9",
      inputPrimary: '#81D4FA',
      codeSectionBg: '#2C2D41',
    },
    extend: {},
  },
  plugins: [],
};
