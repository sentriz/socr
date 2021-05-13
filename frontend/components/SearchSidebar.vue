<template>
  <transition-fade>
    <div v-if="media" class="fixed inset-0 z-10 transition-opacity bg-gray-700 bg-opacity-75" />
  </transition-fade>
  <transition-slide>
    <div v-if="media" ref="content" class="fixed inset-y-0 right-0 z-20 w-full max-w-lg">
      <div class="bg-gray-50 h-full p-6 space-y-6 overflow-y-auto">
        <search-sidebar-header :hash="media.hash" />
        <media-background :hash="media.hash" class="box p-3">
          <media-highlight :hash="media.hash" class="mx-auto" />
        </media-background>
        <div v-if="blocks.length" class="box padded font-mono text-sm bg-gray-200">
          <p v-for="(block, i) in blocks" :key="i" :class="{ 'bg-yellow-300': highlightedBlocksIndexes.has(i) }">
            {{ block.body }}
          </p>
        </div>
      </div>
    </div>
  </transition-slide>
</template>

<script setup lang="ts">
import MediaHighlight from './MediaHighlight.vue'
import TransitionFade from './TransitionFade.vue'
import TransitionSlide from './TransitionSlideX.vue'
import MediaBackground from './MediaBackground.vue'
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

// load the media from the network if we can't find it in the store
// (can happen on page reload if we've click an image on the eg. 5th page)
watch(
  () => props.hash,
  (id) => id && store.loadMedia(id),
  { immediate: true },
)

const media = computed(() => store.getMediaByHash(props.hash || ''))
const blocks = computed(() => store.getBlocksByHash(props.hash || ''))
const highlightedBlocksIndexes = computed(() => {
  const hashBlocks = store.getHighlightedBlocksByHash(props.hash || '')
  return new Set(hashBlocks.map((blocks) => blocks.index))
})

const content = ref<HTMLElement>()
onClickOutside(content, () => router.push({ name: 'search' }))
</script>
