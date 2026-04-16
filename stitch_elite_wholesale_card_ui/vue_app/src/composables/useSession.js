import { computed, ref } from "vue";
import { clearSession, getSession, saveSession } from "../services/api";

const sessionState = ref(getSession());

export function useSession() {
  const isAuthenticated = computed(() => Boolean(sessionState.value?.token));

  function setSession(data) {
    sessionState.value = saveSession(data);
  }

  function logout() {
    clearSession();
    sessionState.value = null;
  }

  function refreshSession() {
    sessionState.value = getSession();
  }

  return {
    session: sessionState,
    isAuthenticated,
    setSession,
    logout,
    refreshSession,
  };
}
