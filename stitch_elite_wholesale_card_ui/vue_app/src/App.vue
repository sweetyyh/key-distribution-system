<script setup>
import AppHeader from "./components/AppHeader.vue";
</script>

<template>
  <div class="app-shell">
    <div class="ambient ambient-a"></div>
    <div class="ambient ambient-b"></div>
    <AppHeader />
    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>
