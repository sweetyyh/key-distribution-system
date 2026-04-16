<script setup>
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useSession } from "../composables/useSession";

const router = useRouter();
const route = useRoute();
const { session, isAuthenticated, logout } = useSession();

const isAuthPage = computed(() =>
  route.name === "login" || route.name === "register",
);

function handleLogout() {
  logout();
  router.push({ name: "home" });
}
</script>

<template>
  <header class="app-header">
    <div class="brand-cluster">
      <RouterLink class="brand-lockup" to="/">
        <span class="brand-mark"></span>
        <span>
          <strong>Stitch Elite</strong>
          <small>Wholesale Card Exchange</small>
        </span>
      </RouterLink>
      <nav class="header-nav" v-if="!isAuthPage">
        <RouterLink to="/">Marketplace</RouterLink>
        <RouterLink to="/orders">Orders</RouterLink>
      </nav>
    </div>

    <div class="header-actions">
      <template v-if="isAuthenticated">
        <span class="eyebrow-chip subtle">{{ session?.username }}</span>
        <button class="button button-secondary" type="button" @click="handleLogout">
          Sign Out
        </button>
      </template>
      <template v-else>
        <RouterLink class="button button-secondary" :to="{ name: 'login' }">
          Sign In
        </RouterLink>
        <RouterLink class="button button-primary" :to="{ name: 'register' }">
          Create Buyer Account
        </RouterLink>
      </template>
    </div>
  </header>
</template>
