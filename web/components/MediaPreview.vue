<template>
  <media-background ref="imgParent" v-if="media" :hash="media.hash" class="flex justify-center" :class="{ 'p-2': childNarrow }">
    <media-highlight ref="imgChild" :hash="media.hash" />
  </media-background>
  <loading-spinner v-else class="bg-gray-100" text="processing image" />

  <details v-if="!isVideo && blocks.length" class="box padded bg-gray-100 bg-gray-100 font-mono text-sm">
    <summary class="select-none py-1 text-sm text-gray-500 hover:cursor-pointer">view text</summary>
    <p v-for="(block, i) in blocks" :key="i" class="overflow-x-hidden rounded-lg" :class="{ 'bg-yellow-200/90': highlightedBlocksIndexes.has(i) }">{{ block.body }}</p>
  </details>
  <loading-spinner v-if="!isVideo && !media?.processed" class="bg-gray-100" text="processing text" />
</template>

<script setup lang="ts">
import MediaBackground from './MediaBackground.vue'
import MediaHighlight from './MediaHighlight.vue'
import LoadingSpinner from './LoadingSpinner.vue'
import { MediaType } from '~/request'
import { computed, ref } from 'vue'
import useStore from '~/composables/useStore'
import { useElementSize } from '@vueuse/core'

const props = defineProps<{
  hash: string
}>()

const store = useStore()

const media = computed(() => store.getMediaByHash(props.hash || ''))
const blocks = computed(() => store.getBlocksByHash(props.hash || ''))
const highlightedBlocksIndexes = computed(() => {
  const hashBlocks = store.getHighlightedBlocksByHash(props.hash || '')
  return new Set(hashBlocks.map((blocks) => blocks.index))
})

const isVideo = computed(() => media.value?.type === MediaType.Video)

const imgParent = ref()
const imgParentW = useElementSize(imgParent)
const imgChild = ref()
const imgChildW = useElementSize(imgChild)

const childNarrow = computed(() => imgChildW.width.value < imgParentW.width.value)
</script>
