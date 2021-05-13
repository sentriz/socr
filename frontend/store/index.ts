import { reactive, readonly, InjectionKey } from 'vue'
import { reqSearch, reqMedia, isError, Block } from '../api'
import type { Reponse, Media, Search, PayloadSearch } from '../api'

const mediasLoadState = async (state: State, resp: Media[]) => {
  for (const media of resp) {
    if (media.blocks) {
      state.blocks.set(media.hash, media.blocks)
      delete media.blocks
    }
    if (media.highlighted_blocks) {
      state.highlighted_blocks.set(media.hash, media.highlighted_blocks)
      delete media.highlighted_blocks
    }
    state.medias.set(media.hash, media)
  }
}

export interface State {
  medias: Map<string, Media>
  blocks: Map<string, Block[]>
  highlighted_blocks: Map<string, Block[]>
  toast: string
}

const createStore = () => {
  const state = reactive<State>({
    medias: new Map(),
    blocks: new Map(),
    highlighted_blocks: new Map(),
    toast: '',
  })

  return {
    state: readonly(state),
    async loadMedias(payload: PayloadSearch): Reponse<Search> {
      const resp = await reqSearch(payload)
      if (isError(resp)) return resp
      if (!resp.result.medias?.length) return resp
      mediasLoadState(state, resp.result.medias)
      return resp
    },
    async loadMedia(hash: string): Reponse<Media> {
      const resp = await reqMedia(hash)
      if (isError(resp)) return resp
      if (!resp.result) return resp
      mediasLoadState(state, [resp.result])
      return resp
    },
    getMediaByHash(hash: string) {
      return state.medias.get(hash)
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
