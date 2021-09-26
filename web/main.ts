import { createApp } from 'vue'
import { RouterView } from 'vue-router'
import { registerSW } from 'virtual:pwa-register'

import router from './router'
import store, { storeSymbol } from './store'
import './main.css'

registerSW({
  onRegistered(e) {
    if (!e) return
    e.addEventListener('fetch', () => {})
    console.info('registered service worker')
  },
})

window.onbeforeunload = () => {
  window.scrollTo(0, 0)
}

const app = createApp(RouterView)
app.use(router)
app.provide(storeSymbol, store)
app.mount('#main')
