<template>
  <div
    class="card w-full md:w-7/12 md:max-w-md dark:bg-base-300 dark:shadow-xl mx-auto mt-10 md:mt-20"
  >
    <div class="card-body">
      <div v-if="success">
        Account created with <span class="link">{{ email }}</span>,
        <RouterLink class="link link-primary no-underline" to="/login">Sign In here</RouterLink>.
      </div>

      <Transition>
        <div v-if="!success">
          <div class="card-title justify-center"><span>Register</span></div>
          <p class="mb-5 text-center">Create new account.</p>
          <form @submit.prevent="registerHandler">
            <div>
              <p class="text-error" v-if="errors?.email">
                {{ validationMessage(errors.email) }}
              </p>
              <input
                type="email"
                placeholder="Email"
                :class="{
                  'input-error': errors?.email,
                  'dark:input-error': errors?.email,
                }"
                class="input input-bordered input-accent dark:input-secondary w-full mb-4"
                v-model="email"
              />
            </div>
            <div>
              <p class="text-error" v-if="errors?.firstname">
                {{ validationMessage(errors.firstname) }}
              </p>
              <input
                type="text"
                placeholder="First Name"
                :class="{
                  'input-error': errors?.firstname,
                  'dark:input-error': errors?.firstname,
                }"
                class="input input-bordered input-accent dark:input-secondary w-full mb-4"
                v-model="firstname"
              />
            </div>
            <div>
              <p class="text-error" v-if="errors?.lastname">
                {{ validationMessage(errors.lastname) }}
              </p>
              <input
                type="text"
                placeholder="Last Name"
                :class="{
                  'input-error': errors?.lastname,
                  'dark:input-error': errors?.lastname,
                }"
                class="input input-bordered input-accent dark:input-secondary w-full mb-4"
                v-model="lastname"
              />
            </div>
            <div>
              <p class="text-error" v-if="errors?.password">
                {{ validationMessage(errors.password) }}
              </p>
              <input
                type="password"
                placeholder="Password"
                :class="{
                  'input-error': errors?.password,
                  'dark:input-error': errors?.password,
                }"
                class="input input-bordered input-accent dark:input-secondary w-full mb-4"
                v-model="password"
              />
            </div>
            <div class="flex justify-between items-center mt-6" v-if="!loading">
              <RouterLink class="link-primary" to="/login">
                Have account, Login!
              </RouterLink>
              <button class="btn btn-sm btn-accent px-4 h-auto normal-case">
                <span class="py-3">Register</span>
              </button>
            </div>
            <div v-if="loading">
              <LoadingLogo class="mx-auto" />
            </div>
          </form>
        </div>
      </Transition>
    </div>
  </div>
</template>

<script>
import agents from "@/agent";
import validation from "@/validation";
import LoadingLogo from "@/components/LoadingLogo.vue";

export default {
  components: { LoadingLogo },
  data: () => ({
    email: "",
    firstname: "",
    lastname: "",
    password: "",
    errors: {},
    loading: true,
    success: false,
  }),
  mounted() {
    this.loading = false;
  },
  methods: {
    async registerHandler() {
      this.errors = {}; // reset errors
      this.loading = true;

      agents.Users.register({
        email: this.email,
        firstname: this.firstname,
        lastname: this.lastname,
        password: this.password,
      })
        .then(({ res, raw }) => {
          if (res === null) {
            alert(raw.statusText);
            return;
          }

          if (raw.status === 400) {
            this.assembleErrors(res.errors);
            return;
          }

          if (raw.status === 201) {
            if (this.email === res?.email) {
              this.success = true;
              return;
            }
          }
        })
        .finally(() => {
          this.loading = false;
        });
    },
    assembleErrors(errors) {
      errors.forEach((e) => {
        this.errors[e.field] = e;
      });
    },
    validationMessage(error) {
      return validation.getErrorMessage(error.field, error.tag, error.param);
    },
  },
};
</script>
