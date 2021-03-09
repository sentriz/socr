import { createApp } from 'vue'
import './main.css'

import App from './components/App.vue'
import router from './router'
import store, { storeSymbol } from './store'

window.onbeforeunload = () => {
  window.scrollTo(0, 0)
}

const app = createApp(App)
app.use(router)
app.provide(storeSymbol, store)
app.mount('body')
