<template>
  <div class="container mx-auto px-4">
    <div class="w-full md:w-8/12 mx-auto mt-5 md:mt-36">
      <div class="text-center" v-if="error">
        <h1 class="text-3xl mb-6">Authorization Error</h1>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-10 h-10 inline"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
      </div>
      <div v-if="!error">
        <h1 class="text-center text-3xl mb-6">Heads Up!</h1>
        <p class="text-center">
          An application requests authorization to your Semaphore account.
        </p>
        <p class="text-center">client-id: {{ queries["client_id"] }}</p>

        <div class="flex gap-2 mt-6 justify-center">
          <button class="btn btn-ghost" @click="handleCancel">Cancel</button>
          <button class="btn btn-success" @click="handleAuthorize">
            Authorize
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "AuthorizeView",
  data: () => ({
    queries: "",
    error: "",
  }),
  created() {
    this.queries = this.$route.query;
  },
  methods: {
    handleCancel() {
      this.$router.push({ name: "dashboard" });
    },
    handleAuthorize: async function () {
      // const authStore = useAuthStore();
      let authorizeAuth =
        "/oauth2/authorize" + window.location.search + "&consent=1";

      const formAuth = document.createElement("form");
      formAuth.setAttribute("action", authorizeAuth);
      formAuth.setAttribute("method", "POST");
      formAuth.hidden = true;

      document.querySelector("body").append(formAuth);
      formAuth.submit();
    },
  },
};
</script>

<style scoped></style>
