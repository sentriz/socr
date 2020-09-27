import { createApp } from "vue";
import { createRouter, createWebHashHistory, RouterView } from "vue-router";

import "./main.css";
import Search from "./components/Search.vue";
import Settings from "./components/Settings.vue";
import SearchSidebar from "./components/SearchSidebar.vue";
import Login from "./components/Login.vue";
import Home from "./components/Home.vue";

import { urlSocket } from './api'

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
          redirect: "search",
        },
      ]
    },
    {
      path: "/:catchAll(.*)",
      redirect: "home",
    },
  ],
});

const socket = new WebSocket(`wss://${window.location.host}${urlSocket}`);

const app = createApp(RouterView);
app.use(router);
app.provide("socket", socket)
app.mount("#app");

window.router = router
