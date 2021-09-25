import { createRouter, createWebHistory, NavigationGuard, RouteRecordRaw } from 'vue-router'

import Search from '~/components/Search.vue'
import Settings from '~/components/Settings.vue'
import Importer from '~/components/Importer.vue'
import Login from '~/components/Login.vue'
import Home from '~/components/Home.vue'
import Public from '~/components/Public.vue'
import NotFound from '~/components/NotFound.vue'

import { tokenHas, tokenSet } from '~/request'

const beforeCheckAuth: NavigationGuard = (to, _, next) => {
  if (tokenHas()) return next()
  next({ name: 'login', query: { redirect: to.fullPath } })
}

const beforeLogout: NavigationGuard = (_, __, next) => {
  tokenSet('')
  next({ name: 'login' })
}

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: Login,
  },
  {
    path: '/logout',
    name: 'logout',
    beforeEnter: beforeLogout,
    redirect: '',
  },
  {
    path: '/',
    name: 'home',
    component: Home,
    redirect: 'search',
    beforeEnter: beforeCheckAuth,
    children: [
      {
        path: 'search/:hash?',
        name: 'search',
        component: Search,
      },
      {
        path: 'importer',
        name: 'importer',
        component: Importer,
      },
      {
        path: 'settings',
        name: 'settings',
        component: Settings,
      },
      {
        path: '',
        redirect: { name: 'search' },
      },
    ],
  },
  {
    path: '/i/:hash',
    name: 'public',
    component: Public,
  },
  {
    path: '/not-found',
    name: 'not-found',
    component: NotFound,
  },
  {
    path: '/:catchAll(.*)',
    redirect: { name: 'not-found' },
  },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
