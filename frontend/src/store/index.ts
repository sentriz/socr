import { reactive, readonly, InjectionKey } from "vue";
import { reqSearch, reqScreenshot, isError, PayloadSort, Block } from "../api";
import type { Reponse, Screenshot, Search } from "../api";

const screenshotsLoadState = async (state: State, resp?: Screenshot[]) => {
  if (!resp) return
  for (const screenshot of resp || []) {
    if (screenshot.blocks) {
      state.blocks[screenshot.hash] = screenshot.blocks
      delete screenshot.blocks
    }
    if (screenshot.highlighted_blocks) {
      state.highlighted_blocks[screenshot.hash] = screenshot.highlighted_blocks
      delete screenshot.highlighted_blocks
    }
    state.screenshots[screenshot.hash] = screenshot;
  }
};

export interface State {
  screenshots: {[hash: string]: Screenshot},
  blocks: {[hash: string]: Block[]},
  highlighted_blocks: {[hash: string]: Block[]},
  toast: string,
}

const createStore = () => {
  const state = reactive<State>({
    screenshots: {},
    blocks: {},
    highlighted_blocks: {},
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
    getBlocksByHash(hash: string) {
      return state.blocks[hash] || [];
    },
    getHighlightedBlocksByHash(hash: string) {
      return state.highlighted_blocks[hash] || [];
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