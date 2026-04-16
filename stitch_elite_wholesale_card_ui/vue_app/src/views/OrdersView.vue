<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { useSession } from "../composables/useSession";
import { apiRequest } from "../services/api";
import {
  formatDate,
  formatMoney,
  getPayChannelMeta,
  getStatusMeta,
} from "../services/catalog";

const route = useRoute();
const router = useRouter();
const { isAuthenticated, refreshSession } = useSession();

const searchKeyword = ref("");
const loading = ref(false);
const orderLoading = ref(false);
const pendingPayment = ref(false);
const orders = ref([]);
const orderDetail = ref(null);
const cards = ref([]);
const selectedOrderNo = ref("");
const listError = ref("");
const detailError = ref("");
const cardsError = ref("");

const filteredOrders = computed(() => {
  const value = searchKeyword.value.trim().toLowerCase();
  if (!value) return orders.value;
  return orders.value.filter((order) => {
    return (
      order.order_no.toLowerCase().includes(value) ||
      order.product_name.toLowerCase().includes(value)
    );
  });
});

function toneClass(status) {
  const tone = getStatusMeta(status).tone;
  return {
    "tone-success": tone === "success",
    "tone-warning": tone === "warning",
    "tone-muted": tone === "muted",
  };
}

async function loadOrders() {
  refreshSession();
  if (!isAuthenticated.value) return;

  loading.value = true;
  listError.value = "";
  try {
    const data = await apiRequest("/api/v1/orders?page=1&page_size=20", {
      auth: true,
    });
    orders.value = Array.isArray(data.list) ? data.list : [];

    const requestedOrder = typeof route.query.order_no === "string" ? route.query.order_no : "";
    if (requestedOrder) {
      selectedOrderNo.value = requestedOrder;
    } else if (!selectedOrderNo.value && orders.value.length > 0) {
      selectedOrderNo.value = orders.value[0].order_no;
    }

    if (selectedOrderNo.value) {
      await loadOrderDetail(selectedOrderNo.value);
    }
  } catch (error) {
    listError.value = error.message;
  } finally {
    loading.value = false;
  }
}

async function loadOrderDetail(orderNo) {
  detailError.value = "";
  cardsError.value = "";
  cards.value = [];
  orderLoading.value = true;
  try {
    orderDetail.value = await apiRequest(`/api/v1/orders/${encodeURIComponent(orderNo)}`, {
      auth: true,
    });
    if (orderDetail.value.status === 2) {
      await loadCards();
    }
  } catch (error) {
    detailError.value = error.message;
    orderDetail.value = null;
  } finally {
    orderLoading.value = false;
  }
}

async function loadCards() {
  if (!selectedOrderNo.value) return;
  cardsError.value = "";
  try {
    const data = await apiRequest(
      `/api/v1/orders/${encodeURIComponent(selectedOrderNo.value)}/cards`,
      { auth: true },
    );
    cards.value = Array.isArray(data.cards) ? data.cards : [];
  } catch (error) {
    cardsError.value = error.message;
  }
}

async function selectOrder(orderNo) {
  selectedOrderNo.value = orderNo;
  await router.replace({
    name: "orders",
    query: {
      order_no: orderNo,
    },
  });
  await loadOrderDetail(orderNo);
}

async function simulatePayment() {
  if (!orderDetail.value) return;
  pendingPayment.value = true;
  detailError.value = "";
  try {
    await apiRequest(`/api/v1/pay/callback/${encodeURIComponent(orderDetail.value.pay_channel)}`, {
      method: "POST",
      body: {
        order_no: orderDetail.value.order_no,
      },
    });
    await loadOrders();
  } catch (error) {
    detailError.value = error.message;
  } finally {
    pendingPayment.value = false;
  }
}

async function copyCard(card) {
  const value = card.pin ? `${card.code} | ${card.pin}` : card.code;
  await navigator.clipboard.writeText(value);
}

onMounted(loadOrders);

watch(
  () => route.query.order_no,
  async (nextOrderNo) => {
    if (typeof nextOrderNo === "string" && nextOrderNo && nextOrderNo !== selectedOrderNo.value) {
      selectedOrderNo.value = nextOrderNo;
      await loadOrderDetail(nextOrderNo);
    }
  },
);
</script>

<template>
  <div class="page">
    <section class="surface">
      <div class="surface-content section-heading">
        <div>
          <span class="eyebrow-chip">Operations Console</span>
          <h2 style="margin-top: 0.75rem">Buyer orders, detail, payment simulation, and cards.</h2>
          <p class="section-subtitle">
            This view preserves the current mock payment flow exposed by the backend and wraps it
            in a higher-end operational interface.
          </p>
        </div>
        <label class="search-shell" style="min-width: 280px">
          <input v-model="searchKeyword" placeholder="Search order or product" type="search" />
        </label>
      </div>
    </section>

    <section v-if="!isAuthenticated" class="surface">
      <div class="surface-content">
        <div class="message-box warning">
          Orders require a buyer token. Sign in first, then the SPA will load the protected order
          APIs against the existing backend.
        </div>
        <div style="margin-top: 1rem">
          <RouterLink
            class="button button-primary"
            :to="{ name: 'login', query: { next: route.fullPath } }"
          >
            Sign In To Continue
          </RouterLink>
        </div>
      </div>
    </section>

    <section v-else class="orders-grid">
      <article class="surface orders-panel">
        <div class="surface-content">
          <div class="section-heading">
            <div>
              <h2>Order list</h2>
              <p class="section-subtitle">{{ filteredOrders.length }} order(s) in the current view.</p>
            </div>
            <button class="button button-secondary" type="button" @click="loadOrders">
              Refresh
            </button>
          </div>

          <div v-if="listError" class="message-box error" style="margin-top: 1rem">{{ listError }}</div>
          <div v-else-if="loading" class="empty-state" style="margin-top: 1rem">Loading orders...</div>
          <div v-else-if="filteredOrders.length === 0" class="empty-state" style="margin-top: 1rem">
            No orders matched the current filter.
          </div>
          <div v-else class="orders-list">
            <button
              v-for="order in filteredOrders"
              :key="order.order_no"
              :class="['order-card', { active: selectedOrderNo === order.order_no }]"
              type="button"
              @click="selectOrder(order.order_no)"
            >
              <div class="order-card__top">
                <div>
                  <div class="order-card__id">{{ order.order_no }}</div>
                  <strong>{{ order.product_name }}</strong>
                </div>
                <span :class="toneClass(order.status)">
                  {{ getStatusMeta(order.status).label }}
                </span>
              </div>
              <p class="subtle-copy">
                Qty {{ order.quantity }} / {{ formatMoney(order.total_amount) }}
              </p>
              <small class="subtle-copy">Created {{ formatDate(order.created_at) }}</small>
            </button>
          </div>
        </div>
      </article>

      <article class="surface detail-panel">
        <div class="surface-content">
          <div v-if="orderLoading" class="empty-state">Loading order detail...</div>
          <template v-else-if="orderDetail">
            <div class="detail-toolbar">
              <div>
                <h2>Order detail</h2>
                <p class="section-subtitle">{{ orderDetail.order_no }}</p>
              </div>
              <div style="display: flex; gap: 0.75rem; flex-wrap: wrap">
                <button
                  v-if="orderDetail.status === 0"
                  class="button button-primary"
                  :disabled="pendingPayment"
                  type="button"
                  @click="simulatePayment"
                >
                  {{ pendingPayment ? "Processing..." : "Simulate Payment" }}
                </button>
                <button
                  v-if="orderDetail.status === 2"
                  class="button button-secondary"
                  type="button"
                  @click="loadCards"
                >
                  Reload Cards
                </button>
              </div>
            </div>

            <div v-if="detailError" class="message-box error" style="margin-top: 1rem">{{ detailError }}</div>

            <div class="detail-grid" style="margin-top: 1rem">
              <div class="detail-cell">
                <span>Status</span>
                <strong :class="toneClass(orderDetail.status)">
                  {{ getStatusMeta(orderDetail.status).label }}
                </strong>
              </div>
              <div class="detail-cell">
                <span>Payment</span>
                <strong>{{ getPayChannelMeta(orderDetail.pay_channel).label }}</strong>
              </div>
              <div class="detail-cell">
                <span>Total</span>
                <strong>{{ formatMoney(orderDetail.total_amount) }}</strong>
              </div>
              <div class="detail-cell">
                <span>Product</span>
                <strong>{{ orderDetail.product_name }}</strong>
              </div>
              <div class="detail-cell">
                <span>Quantity</span>
                <strong>{{ orderDetail.quantity }}</strong>
              </div>
              <div class="detail-cell">
                <span>Unit Price</span>
                <strong>{{ formatMoney(orderDetail.unit_price) }}</strong>
              </div>
              <div class="detail-cell">
                <span>Created At</span>
                <strong>{{ formatDate(orderDetail.created_at) }}</strong>
              </div>
              <div class="detail-cell">
                <span>Expires At</span>
                <strong>{{ formatDate(orderDetail.expires_at) }}</strong>
              </div>
              <div class="detail-cell">
                <span>Paid At</span>
                <strong>{{ formatDate(orderDetail.paid_at) }}</strong>
              </div>
            </div>

            <div style="margin-top: 1.4rem">
              <h3 style="margin-bottom: 0.8rem">Delivered cards</h3>
              <div v-if="cardsError" class="message-box error">{{ cardsError }}</div>
              <div v-else-if="orderDetail.status !== 2" class="message-box warning">
                This order is still pending. Use the mock payment callback, then reload cards.
              </div>
              <div v-else-if="cards.length === 0" class="empty-state">
                No cards returned yet.
              </div>
              <div v-else class="cards-list">
                <div v-for="card in cards" :key="`${card.index}-${card.code}`" class="card-code">
                  <div class="card-code__top">
                    <div>
                      <small>Card #{{ card.index }}</small>
                      <div class="code-line">{{ card.code }}</div>
                      <div class="code-line">{{ card.pin || "-" }}</div>
                    </div>
                    <button class="button button-secondary" type="button" @click="copyCard(card)">
                      Copy
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </template>
          <div v-else class="empty-state">
            Select an order from the left side to inspect its full lifecycle.
          </div>
        </div>
      </article>
    </section>
  </div>
</template>
