(function () {
  "use strict";

  var SESSION_KEY = "glacier.portal.session";
  var CHECKOUT_KEY = "glacier.portal.checkout";
  var ORDER_PAGE = "/portal/order_web/order_web.html";
  var LOGIN_PAGE = "/portal/log_in_web/log_in_web.html";
  var INDEX_PAGE = "/portal/index_web/index_web.html";

  function safeJsonParse(value, fallback) {
    if (!value) {
      return fallback;
    }
    try {
      return JSON.parse(value);
    } catch (error) {
      return fallback;
    }
  }

  function normalizePath(path) {
    if (!path || typeof path !== "string") {
      return INDEX_PAGE;
    }
    if (!path.startsWith("/portal/")) {
      return INDEX_PAGE;
    }
    return path;
  }

  function getNextPath(fallback) {
    var params = new URLSearchParams(window.location.search);
    return normalizePath(params.get("next") || fallback || INDEX_PAGE);
  }

  function buildLoginUrl(nextPath) {
    return LOGIN_PAGE + "?next=" + encodeURIComponent(normalizePath(nextPath));
  }

  function saveSession(data) {
    var session = {
      token: data.token,
      userId: data.user_id,
      username: data.username,
      role: data.role,
    };
    window.localStorage.setItem(SESSION_KEY, JSON.stringify(session));
    return session;
  }

  function getSession() {
    var session = safeJsonParse(window.localStorage.getItem(SESSION_KEY), null);
    if (!session || !session.token) {
      return null;
    }
    return session;
  }

  function clearSession() {
    window.localStorage.removeItem(SESSION_KEY);
  }

  function saveCheckoutDraft(data) {
    window.localStorage.setItem(CHECKOUT_KEY, JSON.stringify(data || {}));
  }

  function getCheckoutDraft() {
    return safeJsonParse(window.localStorage.getItem(CHECKOUT_KEY), {});
  }

  async function request(path, options) {
    var settings = options || {};
    var headers = new Headers(settings.headers || {});
    headers.set("Accept", "application/json");

    if (!(settings.body instanceof FormData)) {
      headers.set("Content-Type", "application/json");
    }

    if (settings.auth) {
      var session = getSession();
      if (!session) {
        var authError = new Error("Please log in first.");
        authError.status = 401;
        authError.code = 40100;
        throw authError;
      }
      headers.set("Authorization", "Bearer " + session.token);
    }

    var response = await fetch(path, {
      method: settings.method || "GET",
      headers: headers,
      body:
        settings.body instanceof FormData
          ? settings.body
          : settings.body
            ? JSON.stringify(settings.body)
            : undefined,
    });

    var payload = null;
    try {
      payload = await response.json();
    } catch (error) {
      payload = null;
    }

    if (!response.ok || !payload || payload.code !== 0) {
      var message =
        (payload && payload.message) || "Request failed with status " + response.status;
      var requestError = new Error(message);
      requestError.status = response.status;
      requestError.code = payload ? payload.code : -1;
      requestError.payload = payload;
      if (response.status === 401 || response.status === 403) {
        clearSession();
      }
      throw requestError;
    }

    return payload.data;
  }

  function requireLogin(nextPath) {
    window.location.href = buildLoginUrl(nextPath || window.location.pathname + window.location.search);
  }

  function formatMoney(value) {
    var number = Number(value || 0);
    if (!Number.isFinite(number)) {
      return "$0.00";
    }
    return "$" + number.toFixed(2);
  }

  function formatDate(value) {
    if (!value) {
      return "-";
    }
    var date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return value;
    }
    return new Intl.DateTimeFormat("zh-CN", {
      dateStyle: "medium",
      timeStyle: "short",
    }).format(date);
  }

  function calculateUnitPrice(rules, quantity, fallbackPrice) {
    var safeRules = Array.isArray(rules) ? rules : [];
    for (var index = 0; index < safeRules.length; index += 1) {
      var rule = safeRules[index];
      if (quantity >= rule.min && (!rule.max || quantity <= rule.max)) {
        return Number(rule.price || 0);
      }
    }
    if (safeRules.length > 0) {
      return Number(safeRules[safeRules.length - 1].price || 0);
    }
    return Number(fallbackPrice || 0);
  }

  function getStatusMeta(status) {
    if (status === 2) {
      return {
        label: "Paid",
        className: "bg-emerald-500/10 text-emerald-300 border border-emerald-500/20",
      };
    }
    if (status === 0) {
      return {
        label: "Pending Payment",
        className: "bg-amber-500/10 text-amber-300 border border-amber-500/20",
      };
    }
    return {
      label: "Unknown",
      className: "bg-slate-500/10 text-slate-200 border border-slate-500/20",
    };
  }

  function getPayChannelMeta(channel) {
    var value = (channel || "").toLowerCase();
    if (value === "usdt") {
      return { label: "USDT", icon: "currency_bitcoin" };
    }
    if (value === "alipay") {
      return { label: "Alipay", icon: "account_balance_wallet" };
    }
    if (value === "card") {
      return { label: "Card", icon: "credit_card" };
    }
    return { label: channel || "Unknown", icon: "payments" };
  }

  function productArtwork(productId) {
    var images = {
      1: "https://lh3.googleusercontent.com/aida-public/AB6AXuDhPKqV3ZeFQkFweG_H4zl9u7puTrrb6Y3WqqMV8JD7T9k-LqBlEWWP9b43bevOsHjHmXSPttd_7xEDB8Zmbm_V0KQ3is1zm7yMgsP2PNncY7PTzSWzmSuZPds6hEA9T6dcmH6Pa7t8_gPRnKn2Jpk13YYvmpBdGFetYpBGe9Ej5l31HAvAbAbCFSH9X_tY2yPlSEb2AddSUwMiNqUBp0iDSbNmV8l38BgUnvWRQWVk2ezTp1507juInFxK0Xdnm7ONObkd6iKVN_s",
      2: "https://lh3.googleusercontent.com/aida-public/AB6AXuDA8eyKKranRD859p_1DjYq7vTQ5mwl05mgN3u5jFF6SKJfxhU4RsL3rPyVn794FLvAtY3hMmMsdRHOzsuI6RGinbEGzYA0oD2QmpcS0cC95xfNyjh6L6y96L1anzGvAjavhKZ5nUKdrlC21pAoQJDJbV6bIfcdQItkONamDsRUwqji1brq1V6EQiSBgrKCugInsDq8K31GHuO497qAVhqsbatGx4PjNkHj5arvi29w3YBXIQQ2ulDT3xOc_ZSDTQi_4fzMK669lHs",
      3: "https://lh3.googleusercontent.com/aida-public/AB6AXuDA8eyKKranRD859p_1DjYq7vTQ5mwl05mgN3u5jFF6SKJfxhU4RsL3rPyVn794FLvAtY3hMmMsdRHOzsuI6RGinbEGzYA0oD2QmpcS0cC95xfNyjh6L6y96L1anzGvAjavhKZ5nUKdrlC21pAoQJDJbV6bIfcdQItkONamDsRUwqji1brq1V6EQiSBgrKCugInsDq8K31GHuO497qAVhqsbatGx4PjNkHj5arvi29w3YBXIQQ2ulDT3xOc_ZSDTQi_4fzMK669lHs",
      4: "https://lh3.googleusercontent.com/aida-public/AB6AXuD36eWRICgA2PZTkj1G6Dn6NMm3kbiz98ZzYjqCOofwtpg2PvDGF-EeMbvt8Xu8Ji_tJbOha7J52k-NHR1MbLLaTzjZjUhYStO-BTCMvxdltlCoiNRxLctgb71vClHqwK23XefUaDhSciqtRse9KhEdfPH6eGvZ-WlfemVYb3T5mt831_8GJHTkflkpJVZ3OVWn9AAmECv_XcvsVmQz0A-aIhZdBckUjwIi_1osR8kRdzgGkPtMR74izCqTFPm_KIJ2MKI7ASbHUKE",
      5: "https://lh3.googleusercontent.com/aida-public/AB6AXuAy_KAaoWOq4yXEd22HNjaScGxuWW6Mm2A3QMo2S7_I1oBzlfRsAxQjRr_q3bcpHvOJYbHGxIPeD0nVUm-BbLL_H2KjDQ9SeLZMiIejCMYGxjuT0ip4AmhxalxruNSAc8cVKuIp_4yoEflz2MG6S-ZnHJutktlYZHMeq4HYnnthkbbyDU1u_R7LYbwZYmOa-MEJmCXy6C37HO3WaVuOJlux6krgj9xT3bXVRHn4_AjwmyfDuJrdZgcw6XUcnK9d8KDuUuDXnan7AJ4",
      6: "https://lh3.googleusercontent.com/aida-public/AB6AXuANIXUqLFLTvw1npYKE4I11XxI89D_5t-JawPtyqLiSgpX_7ruOvtOfcj0KfVe1-C7SJ3E3z3iQhCy1UyFKSFMlRn8juihCS1YlrtFG3eDNUjT_qnmZDkbnZGWyKAVQoU7hzjTej1gOYADHelGXhM5abkvZazECt0Dj5awvaCf0Blyp3YGYKxO7ZieQAW7DxQkwiBvrAd60eek15TcJOneHnOKSRP7bJb1ED93gavjVeR2-shBN0Tik8RfigkPJ7YghRv8p4U48bUY",
    };
    return images[productId] || images[1];
  }

  function escapeHtml(value) {
    return String(value || "")
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll("\"", "&quot;")
      .replaceAll("'", "&#39;");
  }

  window.PortalApp = {
    ORDER_PAGE: ORDER_PAGE,
    LOGIN_PAGE: LOGIN_PAGE,
    INDEX_PAGE: INDEX_PAGE,
    api: request,
    buildLoginUrl: buildLoginUrl,
    calculateUnitPrice: calculateUnitPrice,
    clearSession: clearSession,
    escapeHtml: escapeHtml,
    formatDate: formatDate,
    formatMoney: formatMoney,
    getCheckoutDraft: getCheckoutDraft,
    getNextPath: getNextPath,
    getPayChannelMeta: getPayChannelMeta,
    getSession: getSession,
    getStatusMeta: getStatusMeta,
    productArtwork: productArtwork,
    requireLogin: requireLogin,
    saveCheckoutDraft: saveCheckoutDraft,
    saveSession: saveSession,
  };
})();
