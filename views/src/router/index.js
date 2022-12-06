import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/:menu?",
      name: "dashboard",
      component: () => import("../views/DashView.vue"),
    },
    {
      path: "/login",
      name: "login",
      component: () => import("../views/auth/LoginView.vue"),
      beforeEnter: () => {
        const authStore = useAuthStore();

        if (authStore.isLogged) {
          return { path: "/" };
        }
      },
    },
    {
      path: "/o/oauth/authorize",
      name: "oauth:authorize",
      component: () => import("../views/oauth/AuthorizeView.vue"),
    },
  ],
});

router.beforeEach((to) => {
  const authStore = useAuthStore();
  if (!authStore.isLogged && to.name !== "login") {
    return { name: "login", query: to.query };
  }
});

export default router;
