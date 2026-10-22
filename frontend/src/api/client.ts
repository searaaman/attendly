import axios from "axios";

const TOKEN_KEY = "attendly_token";
const USER_KEY = "attendly_user";

// In dev, Vite proxies /api to the Go server; in Docker, nginx does the same.
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? "/api",
});

api.interceptors.request.use((config) => {
  const token = getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// An expired or invalid token sends the user back to the login page.
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401 && getToken()) {
      clearSession();
      window.location.assign(`${import.meta.env.BASE_URL}login`);
    }
    return Promise.reject(error);
  },
);

export function getToken(): string | null {
  try {
    return localStorage.getItem(TOKEN_KEY);
  } catch {
    return null;
  }
}

export function getUserName(): string {
  try {
    return JSON.parse(localStorage.getItem(USER_KEY) ?? "{}").name ?? "";
  } catch {
    return "";
  }
}

export function saveSession(token: string, user: { name: string }) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}

// Pulls the server's {"error": "..."} message out of a failed request.
export function errorMessage(err: unknown): string {
  if (axios.isAxiosError(err)) {
    const serverMessage = err.response?.data?.error;
    if (typeof serverMessage === "string") return serverMessage;
    if (!err.response) return "Can't reach the server. Is the API running?";
    return `Request failed (${err.response.status})`;
  }
  return "Something went wrong";
}

export default api;
