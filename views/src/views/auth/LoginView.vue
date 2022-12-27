<template>
  <div
    class="card w-full md:w-7/12 md:max-w-md dark:bg-base-300 dark:shadow-xl mx-auto mt-10 md:mt-20"
  >
    <div class="card-body">
      <div class="card-title justify-center"><span>Login</span></div>
      <p class="mb-5 text-center">Login to access your account.</p>

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
          class="input input-bordered input-accent dark:input-secondary w-full mb-4"
          v-model="login"
        />
        <input
          type="password"
          placeholder="Password"
          class="input input-bordered input-accent dark:input-secondary w-full mb-4"
          v-model="password"
        />
        <div class="flex justify-between items-center mt-6">
          <RouterLink class="link-primary" to="/register">
            Create account
          </RouterLink>
          <button class="btn btn-sm btn-accent px-4 h-auto normal-case">
            <span class="py-3">Login</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { authNow } from "@/utils/auth";

export default {
  data: () => ({
    login: "",
    password: "",
    error: "",
  }),

  methods: {
    async loginHandler() {
      this.error = "";
      const res = await authNow(this.login, this.password);

      if (!res.success) {
        this.error = res.error;
        return;
      }

      // Check if redirected from oauth/authorize
      const query = this.$route.query;
      if (query.from === "oauth_authorize") {
        // now throw the user back to oauth
        window.location = `${window.location.origin}/oauth2/authorize${window.location.search}`;
        return;
      }

      // reload on login success
      window.location.reload();
    },
  },
};
</script>
