import { createRouter, createWebHistory, NavigationGuard } from "vue-router";

import Search from "../components/Search.vue";
import Settings from "../components/Settings.vue";
import Login from "../components/Login.vue";
import Home from "../components/Home.vue";
import Public from "../components/Public.vue";
import NotFound from "../components/NotFound.vue";

import { tokenHas, tokenSet } from "../api";

const beforeCheckAuth: NavigationGuard = (to, _, next) => {
    if (tokenHas()) return next();
    next({ name: "login", query: { redirect: to.fullPath } });
};

const beforeLogout: NavigationGuard = (_, __, next) => {
    tokenSet("");
    next({name: "login" });
};

export default createRouter({
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
            redirect: ""
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