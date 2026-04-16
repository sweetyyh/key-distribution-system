<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import ProductCard from "../components/ProductCard.vue";
import { apiRequest } from "../services/api";
import { formatMoney } from "../services/catalog";

const loading = ref(true);
const errorMessage = ref("");
const products = ref([]);
const keyword = ref("");

const filteredProducts = computed(() => {
  const value = keyword.value.trim().toLowerCase();
  if (!value) return products.value;
  return products.value.filter((product) => {
    return (
      product.name.toLowerCase().includes(value) ||
      product.description.toLowerCase().includes(value)
    );
  });
});

const totalStock = computed(() =>
  products.value.reduce((sum, product) => sum + Number(product.stock || 0), 0),
);

const bestStartingPrice = computed(() => {
  const prices = products.value
    .map((product) => Number(product.price || 0))
    .filter((price) => Number.isFinite(price) && price > 0);
  return prices.length ? formatMoney(Math.min(...prices)) : "-";
});

const featuredProduct = computed(() => filteredProducts.value[0] || null);

async function loadProducts() {
  loading.value = true;
  errorMessage.value = "";
  try {
    const data = await apiRequest("/api/v1/products");
    products.value = Array.isArray(data.list) ? data.list : [];
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    loading.value = false;
  }
}

onMounted(loadProducts);
</script>

<template>
  <div class="page">
    <section class="hero-grid">
      <article class="surface">
        <div class="surface-content hero-copy">
          <span class="eyebrow-chip">Vue 3 Buyer Portal</span>
          <h1>Luxury execution for wholesale digital inventory.</h1>
          <p>
            The functional flows stay aligned with the current Go backend. The UI is rebuilt as a
            Vue 3 single-page application with a darker editorial tone, stronger hierarchy, and
            cleaner operational flow.
          </p>

          <div class="hero-actions">
            <RouterLink class="button button-primary" to="/orders">Review Orders</RouterLink>
            <RouterLink class="button button-secondary" to="/register">
              Create Buyer Account
            </RouterLink>
          </div>

          <div class="metric-grid">
            <div class="metric-card">
              <span>Live Products</span>
              <strong>{{ loading ? "-" : products.length }}</strong>
            </div>
            <div class="metric-card">
              <span>Total Stock</span>
              <strong>{{ loading ? "-" : totalStock.toLocaleString("en-US") }}</strong>
            </div>
            <div class="metric-card">
              <span>Best Starting Price</span>
              <strong>{{ loading ? "-" : bestStartingPrice }}</strong>
            </div>
          </div>
        </div>
      </article>

      <aside class="showcase-stack">
        <article v-if="featuredProduct" class="surface floating-card">
          <div class="surface-content">
            <span class="eyebrow-chip">Featured Availability</span>
            <img
              class="floating-card__image"
              src="https://images.unsplash.com/photo-1516321318423-f06f85e504b3?auto=format&fit=crop&w=1200&q=80"
              alt="Console dashboard"
            />
            <div class="floating-card__body">
              <h3>{{ featuredProduct.name }}</h3>
              <p>{{ featuredProduct.description }}</p>
            </div>
            <div class="floating-card__footer">
              <div>
                <span>Instant checkout</span>
                <strong>{{ formatMoney(featuredProduct.price) }}</strong>
              </div>
              <RouterLink
                class="button button-secondary"
                :to="{ name: 'product', params: { id: featuredProduct.id } }"
              >
                Inspect
              </RouterLink>
            </div>
          </div>
        </article>
      </aside>
    </section>

    <section class="page">
      <div class="section-heading">
        <div>
          <h2>Marketplace catalog</h2>
          <p class="section-subtitle">
            Search, compare tiers, and jump into product-level checkout without leaving the SPA.
          </p>
        </div>
        <div class="filters">
          <label class="search-shell">
            <input v-model="keyword" placeholder="Search product or description" type="search" />
          </label>
        </div>
      </div>

      <div v-if="errorMessage" class="message-box error">{{ errorMessage }}</div>
      <div v-else-if="loading" class="empty-state">Loading catalog...</div>
      <div v-else-if="filteredProducts.length === 0" class="empty-state">
        No products matched the current filter.
      </div>
      <div v-else class="product-list">
        <ProductCard
          v-for="product in filteredProducts"
          :key="product.id"
          :product="product"
        />
      </div>
    </section>
  </div>
</template>
