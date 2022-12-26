import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";
import { validateLogin } from "./utils/auth";
import {useThemeStore} from "@/stores/theme";

const app = createApp(App);

app.use(createPinia());

const themeStore = useThemeStore();
themeStore.init();

validateLogin()
  .then((isLogged) => {
    if (!isLogged) {
      router.push({ path: "/login", force: true });
    }
  })
  .finally(() => {
    app.use(router);
    app.mount("#app");
  });
