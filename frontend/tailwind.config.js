/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    colors: {
      primaryFill: "#161616",
      borderFill: "#343434",
      textBlockFill: "#1E1E20",
      primaryGeneral: "#3D3DAB",
      secondaryText: "#8D8D98",
      icon: "#BEBEC3",
      red: "#D34242",
      blueTextPath: "#81D4FA",
      blind: "#5E5E5E",
    },
    extend: {},
  },
  plugins: [require("@headlessui/tailwindcss")],
};
