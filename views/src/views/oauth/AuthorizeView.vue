<template>
  <div class="container mx-auto px-4">
    <div class="w-full md:w-8/12 mx-auto mt-5 md:mt-36">
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
</template>

<script>
export default {
  name: "AuthorizeView",
  data: () => ({
    queries: "",
  }),
  created() {
    this.queries = this.$route.query;
  },
  methods: {
    handleCancel() {
      this.$router.push({ name: "dashboard" });
    },
    handleAuthorize() {},
  },
};
</script>

<style scoped></style>
