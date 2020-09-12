import { createApp } from "vue";

import './main.css'
import App from "./components/App.vue";
import Search from "./components/Search.vue";
import Settings from "./components/Settings.vue";

import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      "path": "/search",
      name: "search",
      component: Search
    },
    {
      "path": "/settings",
      name: "settings",
      component: Settings
    },
    {
      path: "/:catchAll(.*)",
      redirect: { name: "search" },
    },
  ],
})


const app = createApp(App);
app.use(router)
app.mount("#app");
