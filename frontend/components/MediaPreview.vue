<template>
  <media-background v-if="media" :hash="media.hash" class="box flex justify-center">
    <media-highlight :hash="media.hash" />
  </media-background>
  <loading-spinner v-else class="bg-gray-100" text="processing image" />

  <div v-if="blocks.length" class="box padded font-mono text-sm bg-gray-200">
    <p
      v-for="(block, i) in blocks"
      :key="i"
      class="rounded-lg"
      :class="{ 'bg-yellow-200/90': highlightedBlocksIndexes.has(i) }"
    >
      {{ block.body }}
    </p>
  </div>
  <loading-spinner v-else class="bg-gray-100" text="processing text" />
</template>

<script setup lang="ts">
import MediaBackground from './MediaBackground.vue'
import MediaHighlight from './MediaHighlight.vue'
import LoadingSpinner from './LoadingSpinner.vue'
import { computed } from 'vue'
import useStore from '../composables/useStore'

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
