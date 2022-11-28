import { Base64 } from "js-base64";
import { useAuthStore } from "../stores/auth";

export async function authCheck(email, password) {
  const cred = { email, password };

  const res = await fetch(`/api/login?check=1`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(cred),
  });

  if (res.status == 200) {
    return true;
  }

  return false;
}

export function authNow(email, password) {
  const formEl = document.createElement("form");
  formEl.setAttribute("action", "/auth/login");
  formEl.setAttribute("method", "post");

  const emailEl = document.createElement("input");
  const passwdEl = document.createElement("password");
  emailEl.setAttribute("type", "hidden");
  passwdEl.setAttribute("type", "hidden");

  emailEl.value = email;
  passwdEl.value = password;

  document.body.append(formEl);
  formEl.submit();
}

export async function authRenew() {
  const res = await fetch(`/api/renew`, {
    method: "POST",
  });

  return res;
}

export async function validateLogin() {
  const res = await authRenew();

  if (res.status != 201) {
    return false;
  }

  return res
    .json()
    .then((tokenPair) => {
      const at = tokenPair.access_token;
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

  const data = JSON.parse(Base64.decode(parts[1]));
  return data;
}
