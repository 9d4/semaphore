import { defineStore } from "pinia";

const AVAILABLE_THEMES = ["smph", "dark"];

const getOsColorScheme = () => {
  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    return "dark";
  }
  return "light";
};

const LocalStorageThemeKey = "semaphore_theme";

export const useThemeStore = defineStore("theme", {
  state: () => ({
    theme: "dark",
  }),
  actions: {
    save() {
      localStorage.setItem(LocalStorageThemeKey, this.theme);
    },
    setHtmlTheme(theme = "") {
      this.theme = (theme === "") ? this.theme : theme;
      document.querySelector("html").setAttribute("data-theme", this.theme);

      if (this.theme === "dark") {
        document.querySelector("html").classList.add("dark");
        return;
      }

      document.querySelector("html").classList.remove("dark");
    },
    setThemeByOS() {
      switch (getOsColorScheme()) {
        case "dark":
          this.theme = "dark";
          break;
        case "light":
          this.theme = "smph";
          break;
      }
      this.setHtmlTheme();
    },
    setThemeFromLocal() {
      const theme = localStorage.getItem(LocalStorageThemeKey);
      if (theme != null && AVAILABLE_THEMES.includes(theme)) {
        this.theme = theme;
      }
      this.setHtmlTheme();
    },
    init() {
      this.setThemeByOS();
      this.setThemeFromLocal();
      this.save();
    },
  },
});
