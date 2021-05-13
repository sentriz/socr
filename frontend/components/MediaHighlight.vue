<template>
  <div class="w-fit relative">
    <img :src="url" />
    <svg
      v-if="media && blocks.length"
      :viewBox="`0 0 ${media.dim_width} ${media.dim_height}`"
      class="absolute inset-0 text-yellow-500 text-opacity-50 fill-current"
    >
      <rect
        v-for="(b, i) in blocks"
        :key="i"
        :x="b.min_x"
        :y="b.min_y"
        :width="b.max_x - b.min_x"
        :height="b.max_y - b.min_y"
      />
    </svg>
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed } from 'vue'
import { urlMedia } from '../api'
import useStore from '../composables/useStore'

const props = defineProps<{
  hash: string
}>()

const store = useStore()

const media = computed(() => store.getMediaByHash(props.hash))
const blocks = computed(() => store.getHighlightedBlocksByHash(props.hash))
const url = computed(() => `${urlMedia}/${media.value?.hash}/raw`)
</script>
