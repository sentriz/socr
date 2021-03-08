import { reactive, readonly, InjectionKey } from "vue";
import { reqSearch, reqScreenshot, isError, PayloadSort } from "../api";
import type { Reponse, Screenshot, Search } from "../api";

const screenshotsLoadState = async (state: State, resp?: Screenshot[]) => {
  if (!resp) return
  for (const screenshot of resp || []) {
    state.screenshots[screenshot.hash] = screenshot;
  }
};

export interface State {
  screenshots: {[hash: string]: Screenshot},
  toast: string,
}

const createStore = () => {
  const state = reactive<State>({
    // map screenshot id -> screenshot
    screenshots: {},
    toast: "",
  });

  return {
    state: readonly(state),
    async loadScreenshots(size: number, from: number, sort: PayloadSort, term: string): Reponse<Search> {
      const resp = await reqSearch({ size, from, sort, term });
      if (isError(resp)) return resp
      screenshotsLoadState(state, resp.result.screenshots);
      return resp;
    },
    async loadScreenshot(hash: string): Reponse<Screenshot> {
      const resp = await reqScreenshot(hash);
      if (isError(resp)) return resp
      screenshotsLoadState(state, [resp.result]);
      return resp;
    },
    getScreenshotByHash(hash: string) {
      return state.screenshots[hash];
    },
    setToast(toast: string) {
      state.toast = toast
      setTimeout(() => state.toast = "", 1500)
    }
  };
};

export default createStore()
export type Store = ReturnType<typeof createStore>
export const storeSymbol: InjectionKey<Store> = Symbol("store");