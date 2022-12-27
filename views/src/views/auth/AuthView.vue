<script setup>
import logo from "@/assets/images/logo.png";
</script>
<template>
  <div class="container mx-auto px-4 font-inter pb-8">
    <div class="flex justify-center mt-10">
      <img class="w-24" :src="logo" />
    </div>
    <Transition name="bounce" mode="out-in">
      <Component :is="view" />
    </Transition>
  </div>
</template>

<script>
import { defineAsyncComponent } from "vue";
import LoadingCenter from "@/components/LoadingCenter.vue";

export default {
  name: "AuthView",
  components: {
    LoginView: defineAsyncComponent({
      loader: () => import("@/views/auth/LoginView.vue"),
      loadingComponent: LoadingCenter,
    }),
    RegisterView: defineAsyncComponent({
      loader: () => import("@/views/auth/RegisterView.vue"),
      loadingComponent: LoadingCenter,
    }),
  },
  data: () => ({
    view: "",
  }),
  created() {
    this.$watch(
      () => this.$route.name,
      () => {
        this.updateView();
      },
      { immediate: true }
    );
  },
  methods: {
    updateView() {
      switch (this.$route.name) {
        case "login":
          this.view = "LoginView";
          break;
        case "register":
          this.view = "RegisterView";
          break;
      }
    },
  },
};
</script>

<style>
.bounce-enter-active {
  animation: bounce-in 0.5s;
}
.bounce-leave-active {
  animation: bounce-in 0.5s reverse;
}
@keyframes bounce-in {
  0% {
    transform: scale(0.5);
  }
  50% {
    transform: scale(1.15);
  }
  100% {
    transform: scale(1);
  }
}
</style>
