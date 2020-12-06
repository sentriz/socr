import { createApp } from "vue";
import { createRouter, createWebHistory, RouterView } from "vue-router";

import "./main.css";
import Search from "./components/Search.vue";
import Settings from "./components/Settings.vue";
import Login from "./components/Login.vue";
import Home from "./components/Home.vue";
import Public from "./components/Public.vue";
import NotFound from "./components/NotFound.vue";

import { tokenHas, tokenSet } from "./api";
import { storeSymbol, createStore } from "./store";

const beforeCheckAuth = (to, from, next) => {
  if (tokenHas()) return next();
  next({ name: "login", query: { redirect: to.fullPath } });
};

const beforeLogout = (to, from, next) => {
  tokenSet("");
  next({ name: "login" });
};

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "login",
      component: Login,
    },
    {
      path: "/logout",
      name: "logout",
      beforeEnter: beforeLogout,
    },
    {
      path: "/",
      name: "home",
      component: Home,
      redirect: "search",
      beforeEnter: beforeCheckAuth,
      children: [
        {
          path: "search/:id?",
          name: "search",
          component: Search,
        },
        {
          path: "settings",
          name: "settings",
          component: Settings,
        },
        {
          path: "",
          redirect: { name: "search" },
        },
      ],
    },
    {
      path: "/i/:id",
      name: "public",
      component: Public,
    },
    {
      path: "/not_found",
      name: "not_found",
      component: NotFound,
    },
    {
      path: "/:catchAll(.*)",
      redirect: { name: "not_found" },
    },
  ],
});

window.onbeforeunload = () => {
  window.scrollTo(0, 0);
};

const app = createApp(RouterView);
const store = createStore()

app.use(router);
app.provide(storeSymbol, store)
app.mount("body");
