import { createApp } from 'vue'
import { RouterView } from 'vue-router'

import router from './router'
import store, { storeSymbol } from './store'
import './main.css'

window.onbeforeunload = () => {
  window.scrollTo(0, 0)
}

const app = createApp(RouterView)
app.use(router)
app.provide(storeSymbol, store)
app.mount('#main')
