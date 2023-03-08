import { createRouter, createWebHistory, NavigationGuard, RouteRecordRaw } from 'vue-router'

import Search from '~/components/Search.vue'
import Settings from '~/components/Settings.vue'
import Importer from '~/components/Importer.vue'
import Login from '~/components/Login.vue'
import Home from '~/components/Home.vue'
import Public from '~/components/Public.vue'
import NotFound from '~/components/NotFound.vue'

import { tokenHas, tokenSet } from '~/request'

export const routes = {
  LOGIN: Symbol(),
  LOGOUT: Symbol(),
  HOME: Symbol(),
  SEARCH: Symbol(),
  IMPORTER: Symbol(),
  SETTINGS: Symbol(),
  PUBLIC: Symbol(),
  NOT_FOUND: Symbol(),
} as const

const beforeCheckAuth: NavigationGuard = (to, _, next) => {
  if (tokenHas()) return next()
  next({ name: routes.LOGIN, query: { redirect: to.fullPath } })
}

const beforeLogout: NavigationGuard = (_, __, next) => {
  tokenSet('')
  next({ name: routes.LOGIN })
}

const records: RouteRecordRaw[] = [
  {
    path: '/login',
    name: routes.LOGIN,
    component: Login,
  },
  {
    path: '/logout',
    name: routes.LOGOUT,
    beforeEnter: beforeLogout,
    redirect: '',
  },
  {
    path: '/',
    name: routes.HOME,
    component: Home,
    redirect: 'search',
    beforeEnter: beforeCheckAuth,
    children: [
      {
        path: 'search/:hash?',
        name: routes.SEARCH,
        component: Search,
      },
      {
        path: 'importer',
        name: routes.IMPORTER,
        component: Importer,
      },
      {
        path: 'settings',
        name: routes.SETTINGS,
        component: Settings,
      },
      {
        path: '',
        name: Symbol(),
        redirect: { name: routes.SEARCH },
      },
    ],
  },
  {
    path: '/i/:hash',
    name: routes.PUBLIC,
    component: Public,
  },
  {
    path: '/not-found',
    name: routes.NOT_FOUND,
    component: NotFound,
  },
  {
    path: '/:catchAll(.*)',
    redirect: { name: routes.NOT_FOUND },
  },
]

export default createRouter({
  history: createWebHistory(),
  routes: records,
})
