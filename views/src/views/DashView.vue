<template>
  <div class="bg-zinc-900">
    <div class="container md:w-7/12 mx-auto px-4 pt-20">
      <div class="tabs text-slate-50 w-full">
        <RouterLink
          to="/"
          class="tab tab-lifted"
          :class="{ 'tab-active': active == '' || active == 'home' }"
          >Home
        </RouterLink>
        <RouterLink
          to="/profile"
          class="tab tab-lifted"
          :class="{ 'tab-active': active == 'profile' }"
          >Profile
        </RouterLink>
        <span class="tab tab-lifted flex-auto pointer-events-none"></span>
      </div>
    </div>
  </div>
  <div class="container md:w-7/12 mx-auto px-4 mt-5">
    <component :is="view" v-bind="currentViewProps" />
  </div>
</template>

<script>
document.getElementsByTagName("html")[0].setAttribute("data-theme", "dracula");
import GreetingTron from "../components/GreetingTron.vue";
import ProfileList from "../components/ProfileList.vue";
import { useAuthStore } from "../stores/auth";

export default {
  setup() {
    const authStore = useAuthStore();
    return { authStore };
  },
  components: { GreetingTron, ProfileList },
  data() {
    return {
      active: "",
      view: null,
      currentViewProps: {},
    };
  },
  created() {
    this.$watch(
      () => this.$route.params,
      () => {
        this.update();
      },
      { immediate: true }
    );
  },
  methods: {
    update() {
      this.handleTabs();
      this.updateView();
    },
    handleTabs() {
      const menu = this.$route.params.menu;
      switch (menu) {
        case "":
        case "profile":
          this.active = menu;
          break;
        default:
          this.$router.push({ path: "/" });
      }

      this.updateView();
    },
    updateView() {
      switch (this.active) {
        case "":
          this.view = "GreetingTron";
          this.currentViewProps = {
            name: this.authStore.jwt.user.firstname,
          };
          break;
        case "profile":
          this.view = "ProfileList";
          this.currentViewProps = {
            claims: this.authStore.jwt,
          };
          break;
        default:
          this.view = null;
      }
    },
  },
};
</script>
