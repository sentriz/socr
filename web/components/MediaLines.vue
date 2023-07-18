<template>
  <details v-if="blocks.length" class="box padded bg-gray-100 font-mono text-sm">
    <summary class="select-none py-1 text-sm text-gray-600 hover:cursor-pointer">view text</summary>
    <p v-for="(block, i) in blocks" :key="i" class="overflow-x-hidden rounded-lg" :class="{ 'bg-yellow-200/90': highlightedBlocksIndexes.has(i) }">{{ block.body }}</p>
  </details>
  <loading-spinner v-else-if="!media?.processed" text="processing text" />
</template>

<script setup lang="ts">
import LoadingSpinner from './LoadingSpinner.vue'
import { computed } from 'vue'
import useStore from '~/composables/useStore'

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
</script>
