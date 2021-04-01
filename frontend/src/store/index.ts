import { reactive, readonly, InjectionKey } from 'vue'
import { reqSearch, reqScreenshot, isError, Block } from '../api'
import type { Reponse, Screenshot, Search, PayloadSearch } from '../api'

const screenshotsLoadState = async (state: State, resp?: Screenshot[]) => {
  if (!resp) return
  for (const screenshot of resp || []) {
    if (screenshot.blocks) {
      state.blocks.set(screenshot.hash, screenshot.blocks)
      delete screenshot.blocks
    }
    if (screenshot.highlighted_blocks) {
      state.highlighted_blocks.set(screenshot.hash, screenshot.highlighted_blocks)
      delete screenshot.highlighted_blocks
    }
    state.screenshots.set(screenshot.hash, screenshot)
  }
}

export interface State {
  screenshots: Map<string, Screenshot>
  blocks: Map<string, Block[]>
  highlighted_blocks: Map<string, Block[]>
  toast: string
}

const createStore = () => {
  const state = reactive<State>({
    screenshots: new Map(),
    blocks: new Map(),
    highlighted_blocks: new Map(),
    toast: '',
  })

  return {
    state: readonly(state),
    async loadScreenshots(payload: PayloadSearch): Reponse<Search> {
      const resp = await reqSearch(payload)
      if (isError(resp)) return resp
      screenshotsLoadState(state, resp.result.screenshots)
      return resp
    },
    async loadScreenshot(hash: string): Reponse<Screenshot> {
      const resp = await reqScreenshot(hash)
      if (isError(resp)) return resp
      screenshotsLoadState(state, [resp.result])
      return resp
    },
    getScreenshotByHash(hash: string) {
      return state.screenshots.get(hash)
    },
    getBlocksByHash(hash: string) {
      return state.blocks.get(hash) || []
    },
    getHighlightedBlocksByHash(hash: string) {
      return state.highlighted_blocks.get(hash) || []
    },
    setToast(toast: string) {
      state.toast = toast
      setTimeout(() => (state.toast = ''), 1500)
    },
  }
}

export default createStore()
export type Store = ReturnType<typeof createStore>
export const storeSymbol: InjectionKey<Store> = Symbol('store')
