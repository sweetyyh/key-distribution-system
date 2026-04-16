const SESSION_KEY = "glacier.portal.session";
const CHECKOUT_KEY = "glacier.portal.checkout";

export function getSession() {
  try {
    const raw = window.localStorage.getItem(SESSION_KEY);
    if (!raw) {
      return null;
    }
    const parsed = JSON.parse(raw);
    return parsed?.token ? parsed : null;
  } catch (error) {
    return null;
  }
}

export function saveSession(data) {
  const session = {
    token: data.token,
    userId: data.user_id,
    username: data.username,
    role: data.role,
  };
  window.localStorage.setItem(SESSION_KEY, JSON.stringify(session));
  return session;
}

export function clearSession() {
  window.localStorage.removeItem(SESSION_KEY);
}

export function getCheckoutDraft() {
  try {
    return JSON.parse(window.localStorage.getItem(CHECKOUT_KEY) || "{}");
  } catch (error) {
    return {};
  }
}

export function saveCheckoutDraft(data) {
  window.localStorage.setItem(CHECKOUT_KEY, JSON.stringify(data || {}));
}

export async function apiRequest(path, options = {}) {
  const headers = new Headers(options.headers || {});
  headers.set("Accept", "application/json");

  if (!(options.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }

  if (options.auth) {
    const session = getSession();
    if (!session) {
      const authError = new Error("Please log in first.");
      authError.status = 401;
      throw authError;
    }
    headers.set("Authorization", `Bearer ${session.token}`);
  }

  const response = await fetch(path, {
    method: options.method || "GET",
    headers,
    body:
      options.body instanceof FormData
        ? options.body
        : options.body
          ? JSON.stringify(options.body)
          : undefined,
  });

  let payload = null;
  try {
    payload = await response.json();
  } catch (error) {
    payload = null;
  }

  if (!response.ok || !payload || payload.code !== 0) {
    const requestError = new Error(
      payload?.message || `Request failed with status ${response.status}`,
    );
    requestError.status = response.status;
    requestError.code = payload?.code ?? -1;
    requestError.payload = payload;
    if (response.status === 401 || response.status === 403) {
      clearSession();
    }
    throw requestError;
  }

  return payload.data;
}
