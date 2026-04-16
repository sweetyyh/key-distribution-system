import { createRouter, createWebHashHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import LoginView from "../views/LoginView.vue";
import RegisterView from "../views/RegisterView.vue";
import ProductView from "../views/ProductView.vue";
import OrdersView from "../views/OrdersView.vue";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", name: "home", component: HomeView },
    { path: "/login", name: "login", component: LoginView },
    { path: "/register", name: "register", component: RegisterView },
    { path: "/product/:id?", name: "product", component: ProductView, props: true },
    { path: "/orders", name: "orders", component: OrdersView },
  ],
  scrollBehavior() {
    return { top: 0 };
  },
});

export default router;
