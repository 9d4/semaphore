import superagent from "superagent";
import { useAuthStore } from "@/stores/auth";

const API_ROOT = `/api`;

const responseBody = (res) => ({
  res: res.body,
  raw: res,
});

const authStore = useAuthStore();
const tokenPlugin = (req) => {
  if (authStore.accessToken) {
    req.set("authorization", `Bearer ${authStore.accessToken}`);
  }
};

const requests = {
  del: (url) =>
    superagent.del(`${API_ROOT}${url}`).use(tokenPlugin).then(responseBody),
  get: (url) =>
    superagent.get(`${API_ROOT}${url}`).use(tokenPlugin).then(responseBody),
  put: (url, body) =>
    superagent
      .put(`${API_ROOT}${url}`, body)
      .use(tokenPlugin)
      .then(responseBody),
  post: (url, body) =>
    superagent
      .post(`${API_ROOT}${url}`, body)
      .use(tokenPlugin)
      .then(responseBody),
};

const Users = {
  get: (userID) => requests.get(`/users/${userID}/profile`),
};

const agents = {
  Users,
};

export default agents;
