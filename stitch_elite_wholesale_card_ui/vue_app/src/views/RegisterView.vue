<script setup>
import { computed, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { useSession } from "../composables/useSession";
import { apiRequest, saveCheckoutDraft } from "../services/api";
import { normalizeAppRoute } from "../services/navigation";

const router = useRouter();
const route = useRoute();
const { setSession } = useSession();

const form = ref({
  username: "",
  email: "",
  password: "",
  confirmPassword: "",
});

const pending = ref(false);
const errorMessage = ref("");
const successMessage = ref("");

const nextRoute = computed(() => {
  const value = route.query.next;
  return typeof value === "string" ? normalizeAppRoute(value) : "/";
});

async function submit() {
  errorMessage.value = "";
  successMessage.value = "";

  if (form.value.username.trim().length < 3) {
    errorMessage.value = "Username must be at least 3 characters.";
    return;
  }
  if (form.value.password.length < 6) {
    errorMessage.value = "Password must be at least 6 characters.";
    return;
  }
  if (form.value.password !== form.value.confirmPassword) {
    errorMessage.value = "Passwords do not match.";
    return;
  }

  pending.value = true;
  try {
    const data = await apiRequest("/api/v1/auth/register", {
      method: "POST",
      body: {
        username: form.value.username.trim(),
        email: form.value.email.trim(),
        password: form.value.password,
      },
    });
    saveCheckoutDraft({
      email: form.value.email.trim(),
      name: form.value.username.trim(),
    });
    setSession(data);
    successMessage.value = "Registration successful. Redirecting...";
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
        <span class="eyebrow-chip">Buyer Onboarding</span>
        <h1>Open a buyer account without changing the existing backend contract.</h1>
        <p>
          Registration still sends <code>username</code>, <code>email</code>, and
          <code>password</code> to <code>/api/v1/auth/register</code>, then stores the buyer token
          returned by the Go API.
        </p>

        <div class="auth-points">
          <div class="auth-point">
            <span>Validation</span>
            <strong>Frontend checks stay aligned with backend minimums</strong>
          </div>
          <div class="auth-point">
            <span>Checkout continuity</span>
            <strong>Email and buyer name are preserved for ordering</strong>
          </div>
          <div class="auth-point">
            <span>Routing</span>
            <strong>Returns to the route that triggered sign-up</strong>
          </div>
        </div>
      </div>
    </section>

    <section class="surface auth-panel">
      <div class="surface-content">
        <div class="section-heading">
          <div>
            <h2>Create Buyer Account</h2>
            <p class="section-subtitle">Provision a new buyer identity for the portal.</p>
          </div>
          <RouterLink class="inline-link" :to="{ name: 'login', query: { next: nextRoute } }">
            Already registered?
          </RouterLink>
        </div>

        <form class="form-grid" @submit.prevent="submit">
          <label class="field-stack">
            <span>Username</span>
            <input v-model="form.username" autocomplete="username" required type="text" />
          </label>

          <label class="field-stack">
            <span>Email</span>
            <input v-model="form.email" autocomplete="email" required type="email" />
          </label>

          <div class="form-grid two-up">
            <label class="field-stack">
              <span>Password</span>
              <input
                v-model="form.password"
                autocomplete="new-password"
                required
                type="password"
              />
            </label>
            <label class="field-stack">
              <span>Confirm Password</span>
              <input
                v-model="form.confirmPassword"
                autocomplete="new-password"
                required
                type="password"
              />
            </label>
          </div>

          <div v-if="errorMessage" class="message-box error">{{ errorMessage }}</div>
          <div v-if="successMessage" class="message-box success">{{ successMessage }}</div>

          <div class="helper-row">
            <small class="subtle-copy">
              This creates a buyer user only. Admin routes remain unchanged.
            </small>
            <button class="button button-primary" :disabled="pending" type="submit">
              {{ pending ? "Creating..." : "Create Account" }}
            </button>
          </div>
        </form>
      </div>
    </section>
  </div>
</template>
