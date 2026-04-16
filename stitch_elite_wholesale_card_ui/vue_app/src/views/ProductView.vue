<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useSession } from "../composables/useSession";
import { apiRequest, getCheckoutDraft, saveCheckoutDraft } from "../services/api";
import {
  calculateUnitPrice,
  formatMoney,
  getProductArtwork,
} from "../services/catalog";

const route = useRoute();
const router = useRouter();
const { isAuthenticated } = useSession();

const product = ref(null);
const loading = ref(true);
const errorMessage = ref("");
const pending = ref(false);

const orderForm = ref({
  quantity: 1,
  email: "",
  name: "",
  payChannel: "usdt",
});

const productId = computed(() => route.params.id || route.query.id || "");
const outOfStock = computed(() => Number(product.value?.stock || 0) <= 0);

const computedUnitPrice = computed(() => {
  if (!product.value) return 0;
  return calculateUnitPrice(
    product.value.wholesale_rules,
    Number(orderForm.value.quantity || 1),
    product.value.price,
  );
});

const totalPrice = computed(
  () => computedUnitPrice.value * Number(orderForm.value.quantity || 1),
);

function clampQuantity(nextValue) {
  const max = Math.max(1, Number(product.value?.stock || 1));
  const safe = Number(nextValue || 1);
  orderForm.value.quantity = Math.min(Math.max(1, safe), max);
}

async function loadProduct() {
  loading.value = true;
  errorMessage.value = "";
  try {
    product.value = await apiRequest(`/api/v1/products/${encodeURIComponent(productId.value)}`);
    clampQuantity(orderForm.value.quantity);
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    loading.value = false;
  }
}

async function createOrder() {
  if (!isAuthenticated.value) {
    await router.push({
      name: "login",
      query: {
        next: route.fullPath,
      },
    });
    return;
  }

  pending.value = true;
  errorMessage.value = "";
  try {
    saveCheckoutDraft({
      email: orderForm.value.email.trim(),
      name: orderForm.value.name.trim(),
    });
    const data = await apiRequest("/api/v1/orders", {
      method: "POST",
      auth: true,
      body: {
        product_id: Number(product.value.id),
        quantity: Number(orderForm.value.quantity),
        pay_channel: orderForm.value.payChannel,
        email: orderForm.value.email.trim(),
        name: orderForm.value.name.trim(),
      },
    });
    await router.push({
      name: "orders",
      query: {
        order_no: data.order_no,
      },
    });
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    pending.value = false;
  }
}

onMounted(() => {
  const draft = getCheckoutDraft();
  orderForm.value.email = draft.email || "";
  orderForm.value.name = draft.name || "";
  loadProduct();
});

watch(productId, () => {
  if (productId.value) {
    loadProduct();
  }
});
</script>

<template>
  <div class="page product-grid">
    <section class="surface product-hero">
      <template v-if="product">
        <img :src="getProductArtwork(product.id)" :alt="product.name" class="product-hero__image" />
        <div class="product-hero__body">
          <span class="eyebrow-chip">{{ outOfStock ? "Out Of Stock" : "Live Product" }}</span>
          <h1>{{ product.name }}</h1>
          <p>{{ product.description }}</p>

          <div class="summary-grid">
            <div class="stat-block">
              <span>Starting Price</span>
              <strong>{{ formatMoney(product.price) }}</strong>
            </div>
            <div class="stat-block">
              <span>Available Stock</span>
              <strong>{{ Number(product.stock || 0).toLocaleString("en-US") }}</strong>
            </div>
            <div class="stat-block">
              <span>Current Tier</span>
              <strong>{{ `${orderForm.quantity}+` }}</strong>
            </div>
          </div>

          <div class="surface" style="margin-top: 1.5rem">
            <div class="surface-content">
              <h2 style="margin-top: 0">Wholesale tiers</h2>
              <table class="tier-table">
                <thead>
                  <tr>
                    <th>Quantity Range</th>
                    <th>Unit Price</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="rule in product.wholesale_rules" :key="`${rule.min}-${rule.max || 0}`">
                    <td>{{ rule.max ? `${rule.min} - ${rule.max}` : `${rule.min}+` }}</td>
                    <td>{{ formatMoney(rule.price) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </template>
      <div v-else-if="loading" class="surface-content empty-state">Loading product...</div>
      <div v-else class="surface-content message-box error">{{ errorMessage }}</div>
    </section>

    <aside class="surface checkout-panel">
      <div class="surface-content">
        <div class="section-heading">
          <div>
            <h2>Checkout</h2>
            <p class="section-subtitle">Create a buyer order against the current backend.</p>
          </div>
          <span class="eyebrow-chip subtle">{{ isAuthenticated ? "Buyer Session" : "Guest" }}</span>
        </div>

        <div class="form-grid">
          <label class="field-stack">
            <span>Quantity</span>
            <div class="quantity-row">
              <button class="button button-secondary" type="button" @click="clampQuantity(orderForm.quantity - 1)">
                -
              </button>
              <input
                v-model.number="orderForm.quantity"
                min="1"
                type="number"
                @input="clampQuantity(orderForm.quantity)"
              />
              <button class="button button-secondary" type="button" @click="clampQuantity(orderForm.quantity + 1)">
                +
              </button>
            </div>
          </label>

          <label class="field-stack">
            <span>Contact Email</span>
            <input v-model="orderForm.email" required type="email" />
          </label>

          <label class="field-stack">
            <span>Buyer Name</span>
            <input v-model="orderForm.name" type="text" />
          </label>

          <div class="field-stack">
            <span>Payment Channel</span>
            <div class="channel-grid">
              <button
                v-for="channel in ['usdt', 'card', 'alipay']"
                :key="channel"
                :class="['channel-card', { active: orderForm.payChannel === channel }]"
                type="button"
                @click="orderForm.payChannel = channel"
              >
                <strong>{{ channel.toUpperCase() }}</strong>
                <span class="subtle-copy">Backend pay_channel</span>
              </button>
            </div>
          </div>

          <div class="summary-card">
            <div class="summary-row">
              <span>Unit Price</span>
              <strong>{{ formatMoney(computedUnitPrice) }}</strong>
            </div>
            <div class="summary-row">
              <span>Quantity</span>
              <strong>{{ orderForm.quantity }}</strong>
            </div>
            <div class="summary-row total">
              <span>Total</span>
              <strong>{{ formatMoney(totalPrice) }}</strong>
            </div>
          </div>

          <div v-if="!isAuthenticated" class="message-box warning">
            Creating an order requires a buyer login. You can inspect the product without a token.
          </div>
          <div v-if="errorMessage && product" class="message-box error">{{ errorMessage }}</div>

          <button
            class="button button-primary"
            :disabled="pending || outOfStock || !product"
            type="button"
            @click="createOrder"
          >
            {{ pending ? "Creating Order..." : outOfStock ? "Out Of Stock" : "Create Order" }}
          </button>
        </div>
      </div>
    </aside>
  </div>
</template>
