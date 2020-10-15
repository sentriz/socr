import { reactive, provide, inject, readonly } from "vue";
import { reqSearch, reqScreenshot } from "../api";

const screenshotsLoadState = async (state, resp) => {
  for (const hit of resp.hits || []) {
    state.screenshots[hit.id] = hit;
  }
};

export const createStore = () => {
  const state = reactive({
    // map screenshot id -> screenshot
    screenshots: {},
  });

  return {
    state: readonly(state),
    async screenshotsLoad(size, from, sort, term) {
      const resp = await reqSearch({ size, from, sort, term });
      screenshotsLoadState(state, resp);
      return resp;
    },
    async screenshotsLoadID(id) {
      const resp = await reqScreenshot(id);
      screenshotsLoadState(state, resp);
      return resp;
    },
    screenshotByID(id) {
      return state.screenshots[id];
    },
  };
};

export const storeSymbol = Symbol("store");
export const useStore = () => inject(storeSymbol);
export const provideStore = () => provide(storeSymbol, createStore());
