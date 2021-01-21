import { createApp } from "vue";
import "./main.css";

import App from "./components/App.vue";
import router from "./router"
import { storeSymbol, createStore } from "./store";

window.onbeforeunload = () => {
  window.scrollTo(0, 0);
};

const app = createApp(App);
const store = createStore()

app.use(router);
app.provide(storeSymbol, store)
app.mount("body");
