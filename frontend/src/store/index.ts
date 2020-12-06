import { reactive, provide, inject, readonly } from "vue";
import { reqSearch, reqScreenshot, Screenshot, FieldSort, ResponseSearch } from "../api";

const screenshotsLoadState = async (state: State, resp: ResponseSearch<Screenshot>) => {
  for (const hit of resp.hits || []) {
    state.screenshots[hit.id] = hit;
  }
};

interface State {
  screenshots: {[id: string]: Screenshot}
}

export const createStore = () => {
  const state = reactive<State>({
    // map screenshot id -> screenshot
    screenshots: {},
  });

  return {
    state: readonly(state),
    async screenshotsLoad(size: number, from: number, sort: FieldSort[], term: string): Promise<ResponseSearch<Screenshot>> {
      const resp = await reqSearch({ size, from, sort, term });
      screenshotsLoadState(state, resp);
      return resp;
    },
    async screenshotsLoadID(id: string) {
      const resp = await reqScreenshot(id);
      screenshotsLoadState(state, resp);
      return resp;
    },
    screenshotByID(id: string) {
      return state.screenshots[id];
    },
  };
};

export const storeSymbol = Symbol("store");
export const useStore = () => inject(storeSymbol);
export const provideStore = () => provide(storeSymbol, createStore());
