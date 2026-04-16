export function normalizeAppRoute(path) {
  if (!path || typeof path !== "string") {
    return "/";
  }

  if (path.startsWith("/portal/app/#")) {
    return path.slice("/portal/app/#".length) || "/";
  }

  if (path === "/portal/" || path === "/portal/index.html") {
    return "/";
  }

  if (path.startsWith("/portal/index_web/index_web.html")) {
    return "/";
  }

  if (path.startsWith("/portal/log_in_web/log_in_web.html")) {
    return "/login";
  }

  if (path.startsWith("/portal/sign_up_web/sign_up_web.html")) {
    return "/register";
  }

  if (path.startsWith("/portal/order_web/order_web.html")) {
    const url = new URL(path, window.location.origin);
    const orderNo = url.searchParams.get("order_no");
    return orderNo ? `/orders?order_no=${encodeURIComponent(orderNo)}` : "/orders";
  }

  if (path.startsWith("/portal/single_web/single_web.html")) {
    const url = new URL(path, window.location.origin);
    const id = url.searchParams.get("id");
    return id ? `/product/${encodeURIComponent(id)}` : "/";
  }

  return path.startsWith("/") ? path : "/";
}
