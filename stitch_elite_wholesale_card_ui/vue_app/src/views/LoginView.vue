<script setup>
import { computed, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { useSession } from "../composables/useSession";
import { apiRequest } from "../services/api";
import { normalizeAppRoute } from "../services/navigation";

const router = useRouter();
const route = useRoute();
const { setSession } = useSession();

const form = ref({
  username: "",
  password: "",
});

const pending = ref(false);
const errorMessage = ref("");
const successMessage = ref("");

const nextRoute = computed(() => {
  const value = route.query.next;
  return typeof value === "string" ? normalizeAppRoute(value) : "/";
});

async function submit() {
  pending.value = true;
  errorMessage.value = "";
  successMessage.value = "";

  try {
    const data = await apiRequest("/api/v1/auth/login", {
      method: "POST",
      body: {
        username: form.value.username.trim(),
        password: form.value.password,
      },
    });
    setSession(data);
    successMessage.value = "Login successful. Redirecting...";
    await router.push(nextRoute.value);
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    pending.value = false;
  }
}
</script>

<template>
  <div class="page auth-grid">
    <section class="surface auth-copy">
      <div class="surface-content">
        <span class="eyebrow-chip">Buyer Authentication</span>
        <h1>Access the wholesale desk with your existing buyer token flow.</h1>
        <p>
          The backend contract stays unchanged: login still submits <code>username</code> and
          <code>password</code> to <code>/api/v1/auth/login</code>. The difference is purely in
          the interaction model and presentation quality.
        </p>

        <div class="auth-points">
          <div class="auth-point">
            <span>Token model</span>
            <strong>Bearer session persisted locally</strong>
          </div>
          <div class="auth-point">
            <span>Navigation</span>
            <strong>Redirects back to the requested route</strong>
          </div>
          <div class="auth-point">
            <span>Scope</span>
            <strong>Buyer portal only, no backend change required</strong>
          </div>
        </div>
      </div>
    </section>

    <section class="surface auth-panel">
      <div class="surface-content">
        <div class="section-heading">
          <div>
            <h2>Sign In</h2>
            <p class="section-subtitle">Use your buyer account to create and manage orders.</p>
          </div>
          <RouterLink class="inline-link" :to="{ name: 'register', query: { next: nextRoute } }">
            Need an account?
          </RouterLink>
        </div>

        <form class="form-grid" @submit.prevent="submit">
          <label class="field-stack">
            <span>Username</span>
            <input v-model="form.username" autocomplete="username" required type="text" />
          </label>

          <label class="field-stack">
            <span>Password</span>
            <input
              v-model="form.password"
              autocomplete="current-password"
              required
              type="password"
            />
          </label>

          <div v-if="errorMessage" class="message-box error">{{ errorMessage }}</div>
          <div v-if="successMessage" class="message-box success">{{ successMessage }}</div>

          <div class="helper-row">
            <small class="subtle-copy">
              Protected requests will attach your buyer bearer token automatically.
            </small>
            <button class="button button-primary" :disabled="pending" type="submit">
              {{ pending ? "Signing In..." : "Sign In" }}
            </button>
          </div>
        </form>
      </div>
    </section>
  </div>
</template>
