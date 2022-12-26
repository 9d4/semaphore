/* eslint-disable no-undef */
module.exports = {
  content: ["./src/**/*.{vue,js,ts}"],
  darkMode: "class",
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        smph: {
          primary: "#ff9933",
          secondary: "#064663",
          accent: "#ffcc33",
          neutral: "#d1d5db",
          "base-100": "#f3f4f6",
          info: "#3ABFF8",
          success: "#36D399",
          warning: "#FBBD23",
          error: "#F87272",
        },
      },
      {
        dark: {
          primary: "#ff9933",
          secondary: "#064663",
          accent: "#FFCC33",
          neutral: "#191D24",
          "base-100": "#1f2937",
          info: "#3ABFF8",
          success: "#36D399",
          warning: "#FBBD23",
          error: "#F87272",
        },
      },
    ],
  },
};
