import { Base64 } from "js-base64";
import { useAuthStore } from "@/stores/auth";

export async function authNow(email, password) {
  const res = await fetch(`/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email,
      password,
    }),
  });

  let ret = { success: true, error: null };

  if (res.status === 200) {
    return ret;
  }

  if (res.status === 401) {
    ret.success = false;
    ret.error = "Credential not found";
    return ret;
  }

  ret.success = false;
  ret.error = res.statusText;
  return ret;
}

export async function authRenew() {
  return await fetch(`/api/renew`, {
    method: "POST",
  });
}

export async function validateLogin() {
  const res = await authRenew();

  if (res.status !== 201) {
    return false;
  }

  return res
    .json()
    .then((tokenPair) => {
      const at = tokenPair["access_token"];
      const jwt = parseToken(at);

      const authStore = useAuthStore();
      authStore.accessToken = at;
      authStore.jwt = jwt;
      return true;
    })
    .catch(() => {
      return false;
    });
}

export function parseToken(token) {
  const parts = token.split(".");

  if (parts.length !== 3) {
    throw new Error("token malformed");
  }

  return JSON.parse(Base64.decode(parts[1]));
}
