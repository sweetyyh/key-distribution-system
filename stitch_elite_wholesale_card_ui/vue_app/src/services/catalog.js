export const productArtworkMap = {
  1: "https://images.unsplash.com/photo-1612817159949-195b6eb9e31a?auto=format&fit=crop&w=1200&q=80",
  2: "https://images.unsplash.com/photo-1625842268584-8f3296236761?auto=format&fit=crop&w=1200&q=80",
  3: "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=1200&q=80",
  4: "https://images.unsplash.com/photo-1574375927938-d5a98e8ffe85?auto=format&fit=crop&w=1200&q=80",
  5: "https://images.unsplash.com/photo-1607853202273-797f1c22a38e?auto=format&fit=crop&w=1200&q=80",
  6: "https://images.unsplash.com/photo-1520607162513-77705c0f0d4a?auto=format&fit=crop&w=1200&q=80",
};

export function getProductArtwork(productId) {
  return productArtworkMap[productId] || productArtworkMap[1];
}

export function calculateUnitPrice(rules, quantity, fallbackPrice) {
  const safeRules = Array.isArray(rules) ? rules : [];
  for (const rule of safeRules) {
    if (quantity >= rule.min && (!rule.max || quantity <= rule.max)) {
      return Number(rule.price || 0);
    }
  }
  if (safeRules.length > 0) {
    return Number(safeRules[safeRules.length - 1].price || 0);
  }
  return Number(fallbackPrice || 0);
}

export function formatMoney(value) {
  const amount = Number(value || 0);
  if (!Number.isFinite(amount)) {
    return "$0.00";
  }
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    minimumFractionDigits: 2,
  }).format(amount);
}

export function formatDate(value) {
  if (!value) {
    return "-";
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return new Intl.DateTimeFormat("zh-CN", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

export function getStatusMeta(status) {
  if (status === 2) {
    return {
      label: "Paid",
      tone: "success",
    };
  }

  if (status === 0) {
    return {
      label: "Pending Payment",
      tone: "warning",
    };
  }

  return {
    label: "Unknown",
    tone: "muted",
  };
}

export function getPayChannelMeta(channel) {
  const value = String(channel || "").toLowerCase();
  if (value === "usdt") {
    return { label: "USDT", icon: "BTC" };
  }
  if (value === "alipay") {
    return { label: "Alipay", icon: "AL" };
  }
  if (value === "card") {
    return { label: "Card", icon: "CC" };
  }
  return { label: channel || "Unknown", icon: "PM" };
}
