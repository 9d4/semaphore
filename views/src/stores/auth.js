import { defineStore } from "pinia";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    accessToken: "",
    jwt: {},
  }),
  getters: {
    isLogged() {
      if (this.jwt.exp !== undefined) {
        return this.jwt.exp > parseInt(Date.now() / 1000);
      }

      return false;
    },
  },
});
