import { inject } from "vue";
import { Store, storeSymbol } from "../store";

export default () => inject<Store>(storeSymbol);