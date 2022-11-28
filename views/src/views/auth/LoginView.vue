<template>
  <div class="container mx-auto px-4">
    <div
      class="card w-full md:w-96 bg-base-300 shadow-xl mx-auto mt-5 md:mt-36"
    >
      <div class="card-body">
        <div class="card-title mb-6">Login</div>
        <div class="alert alert-error shadow-lg mb-3" v-if="error">
          <div>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="stroke-current flex-shrink-0 h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <span>{{ error }}</span>
          </div>
        </div>
        <form @submit.prevent="loginHandler">
          <input
            type="email"
            placeholder="Email"
            class="input w-full mb-4"
            v-model="login"
          />
          <input
            type="password"
            placeholder="Password"
            class="input w-full mb-4"
            v-model="password"
          />
          <div class="flex justify-end">
            <button class="btn btn-accent px-8">Login</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { authCheck, authNow } from "../../utils/auth";
document.getElementsByTagName("html")[0].setAttribute("data-theme", "dark");
export default {
  data: () => ({
    login: "",
    password: "",
    error: "",
  }),

  methods: {
    async loginHandler() {
      this.error = "";
      const ok = await authCheck(this.login, this.password);
      if (ok) {
        return authNow();
      }

      this.error = "Credential does not found";
    },
  },
};
</script>
