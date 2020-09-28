import { createApp } from "vue";
import { createRouter, createWebHashHistory, RouterView } from "vue-router";

import "./main.css";
import Search from "./components/Search.vue";
import Settings from "./components/Settings.vue";
import SearchSidebar from "./components/SearchSidebar.vue";
import Login from "./components/Login.vue";
import Home from "./components/Home.vue";
import Public from "./components/Public.vue";
import NotFound from "./components/NotFound.vue";

import { urlSocket, tokenHas } from './api'

const checkAuth = (to, from, next) => next(
  tokenHas()
    ? undefined
    : { name: "login", query: { redirect: to.fullPath } }
)

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/login",
      name: "login",
      component: Login,
    },
    {
      path: "/",
      name: "home",
      component: Home,
      redirect: "search",
      beforeEnter: checkAuth,
      children: [
        {
          path: "search",
          name: "search",
          component: Search,
          children: [
            {
              path: "result/:id",
              name: "result",
              component: SearchSidebar,
            },
          ],
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
      ]
    },
    {
      path: "/i/:id",
      name: "public",
      component: Public,
    },
    {
      path: "/:catchAll(.*)",
      component: NotFound,
    },
  ],
});

const app = createApp(RouterView);
app.use(router);
app.provide("socket", new WebSocket(`wss://${window.location.host}${urlSocket}`));
app.mount("#app");
