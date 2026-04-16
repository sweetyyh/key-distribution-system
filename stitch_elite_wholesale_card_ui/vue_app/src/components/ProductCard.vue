<script setup>
import { computed } from "vue";
import { RouterLink } from "vue-router";
import { formatMoney, getProductArtwork } from "../services/catalog";

const props = defineProps({
  product: {
    type: Object,
    required: true,
  },
});

const bestTier = computed(() => {
  const rules = props.product.wholesale_rules || [];
  return rules.length ? rules[rules.length - 1] : null;
});
</script>

<template>
  <article class="product-card">
    <img :src="getProductArtwork(product.id)" :alt="product.name" class="product-card__image" />
    <div class="product-card__body">
      <div class="product-card__header">
        <div>
          <h3>{{ product.name }}</h3>
          <p>{{ product.description }}</p>
        </div>
        <span class="eyebrow-chip">{{ formatMoney(product.price) }}</span>
      </div>

      <div class="product-card__meta">
        <div>
          <span>Stock</span>
          <strong>{{ Number(product.stock || 0).toLocaleString("en-US") }}</strong>
        </div>
        <div>
          <span>Best Tier</span>
          <strong>{{ bestTier ? `${bestTier.min}+` : "-" }}</strong>
        </div>
      </div>

      <RouterLink
        class="button button-primary product-card__cta"
        :to="{ name: 'product', params: { id: product.id } }"
      >
        Open Product
      </RouterLink>
    </div>
  </article>
</template>
