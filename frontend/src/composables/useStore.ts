import { inject, InjectionKey, provide } from "vue";
import { createStore, Store, storeSymbol } from "../store";

export default () => inject<Store>(storeSymbol);