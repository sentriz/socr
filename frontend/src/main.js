import { createApp } from "vue";
import { createRouter, createWebHashHistory } from "vue-router";

import "./main.css";
import App from "./components/App.vue";
import Search from "./components/Search.vue";
import Settings from "./components/Settings.vue";
import SearchSidebar from "./components/SearchSidebar.vue";

import { urlSocket } from './api'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/search",
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
      path: "/settings",
      name: "settings",
      component: Settings,
    },
    {
      path: "/:catchAll(.*)",
      redirect: { name: "search" },
    },
  ],
});

const socket = new WebSocket(`wss://${window.location.host}${urlSocket}`);
console.log(socket);

const app = createApp(App);
app.use(router);
app.mount("#app");
