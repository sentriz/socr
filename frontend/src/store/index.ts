import { reactive, readonly, InjectionKey } from "vue";
import { reqSearch, reqScreenshot } from "../api";
import type { Screenshot, FieldSort, ResponseSearch } from "../api";

const screenshotsLoadState = async (state: State, resp: ResponseSearch<Screenshot>) => {
  if (!resp) return
  for (const hit of resp.hits || []) {
    state.screenshots[hit.id] = hit;
  }
};

export interface State {
  screenshots: {[id: string]: Screenshot},
  toast: string,
}

export const createStore = () => {
  const state = reactive<State>({
    // map screenshot id -> screenshot
    screenshots: {},
    toast: "",
  });

  return {
    state: readonly(state),
    async loadScreenshots(size: number, from: number, sort: FieldSort[], term: string): Promise<ResponseSearch<Screenshot>> {
      const resp = await reqSearch({ size, from, sort, term });
      screenshotsLoadState(state, resp);
      return resp;
    },
    async loadScreenshot(id: string): Promise<ResponseSearch<Screenshot>> {
      const resp = await reqScreenshot(id);
      screenshotsLoadState(state, resp);
      return resp;
    },
    getScreenshotByID(id: string) {
      return state.screenshots[id];
    },
    setToast(toast: string) {
      state.toast = toast
      setTimeout(() => state.toast = "", 1500)
    }
  };
};

export type Store = ReturnType<typeof createStore>
export const storeSymbol: InjectionKey<Store> = Symbol("store");