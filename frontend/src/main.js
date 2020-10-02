import { createApp, reactive } from "vue";
import { createRouter, createWebHashHistory, RouterView } from "vue-router";

import "./main.css";
import Search from "./components/Search.vue";
import Settings from "./components/Settings.vue";
import SearchSidebar from "./components/SearchSidebar.vue";
import Login from "./components/Login.vue";
import Home from "./components/Home.vue";
import Public from "./components/Public.vue";
import NotFound from "./components/NotFound.vue";

import { tokenHas, tokenSet, tokenGet } from './api'

const beforeCheckAuth = (to, from, next) => {
  if (tokenHas()) return next()
  next({ name: "login", query: { redirect: to.fullPath } })
}

const beforeLogout = (to, from, next) => {
  tokenSet("")
  next({ name: "login" })
}

const router = createRouter({
  history: createWebHashHistory(),
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

const store = reactive({
  // map screenshot id -> screenshot
  screenshots: {}
})

const app = createApp(RouterView);
app.use(router);
app.provide("store", store);
app.mount("#app");
