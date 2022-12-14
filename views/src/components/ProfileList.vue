<template>
  <div>
    <h1 class="text-3xl mb-2">Profile</h1>
    <p class="text-slate-400">Your user profile.</p>

    <form @submit.prevent="null" class="mt-4">
      <div class="form-control w-full mb-6">
        <label class="label">
          <span class="label-text text-gray-300">Email</span>
        </label>
        <input
          type="text"
          class="input input-bordered input-sm w-full"
          disabled="!ro"
          :value="claims.user.email"
        />
      </div>
      <div class="form-control w-full mb-6">
        <label class="label">
          <span class="label-text text-gray-300">First Name</span>
        </label>
        <input
          type="text"
          class="input input-bordered input-sm w-full"
          disabled="!ro"
          :value="claims.user.firstname"
        />
      </div>
      <div class="form-control w-full mb-6">
        <label class="label">
          <span class="label-text text-gray-300">Last Name</span>
        </label>
        <input
          type="text"
          class="input input-bordered input-sm w-full"
          disabled="!ro"
          :value="claims.user.lastname"
        />
      </div>
    </form>
  </div>
</template>

<script>
export default {
  props: ["claims"],
  data: () => ({
    ro: true,
  }),
};
</script>
