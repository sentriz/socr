<template>
  <TransitionFade>
    <div v-if="screenshot" class="fixed inset-0 z-10 transition-opacity bg-gray-700 bg-opacity-75" />
  </TransitionFade>
  <TransitionSlide>
    <div v-if="screenshot" ref="content" class="fixed inset-y-0 right-0 z-20 w-full max-w-lg">
      <div class="h-full p-6 space-y-6 overflow-y-auto bg-gray-50">
        <SearchSidebarHeader :hash="screenshot.hash" />
        <ScreenshotBackground :hash="screenshot.hash" class="box p-3">
          <ScreenshotHighlight :hash="screenshot.hash" class="mx-auto" />
        </ScreenshotBackground>
        <div v-if="blocks.length" class="box padded font-mono text-sm bg-gray-200">
          <p v-for="(block, i) in blocks" :key="i" :class="{ 'bg-yellow-300': highlightedBlocksIndexes.has(i) }">
            {{ block.body }}
          </p>
        </div>
      </div>
    </div>
  </TransitionSlide>
</template>

<script setup lang="ts">
import ScreenshotHighlight from './ScreenshotHighlight.vue'
import TransitionFade from './TransitionFade.vue'
import TransitionSlide from './TransitionSlideX.vue'
import ScreenshotBackground from './ScreenshotBackground.vue'
import SearchSidebarHeader from './SearchSidebarHeader.vue'

import { computed, defineProps, ref, watch } from 'vue'
import useStore from '../composables/useStore'
import { onClickOutside } from '@vueuse/core'
import { useRouter } from 'vue-router'

const props = defineProps<{
  hash?: string
}>()

const store = useStore()
const router = useRouter()

// load the screenshot from the network if we can't find it in the store
// (can happen on page reload if we've click an image on the eg. 5th page)
watch(
  () => props.hash,
  (id) => id && store.loadScreenshot(id),
  { immediate: true },
)

const screenshot = computed(() => store.getScreenshotByHash(props.hash || ''))
const blocks = computed(() => store.getBlocksByHash(props.hash || ''))
const highlightedBlocksIndexes = computed(() => {
  const hashBlocks = store.getHighlightedBlocksByHash(props.hash || '')
  return new Set(hashBlocks.map((blocks) => blocks.index))
})

const content = ref<HTMLElement>()
onClickOutside(content, () => router.push({ name: 'search' }))
</script>
