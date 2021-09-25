<template>
  <div class="w-fit relative">
    <div v-if="thumb && isVideo" class="top-2 right-2 bg-black/80 absolute p-2 rounded-md">
      <video-camera-icon class="h-5 text-white" />
    </div>
    <video v-if="!thumb && isVideo && media" :controls="true" @loadstart="loaded">
      <source :src="url" :type="media.mime" />
    </video>
    <img v-else :src="url" @load="loaded" />

    <svg
      v-if="media && blocks.length"
      :viewBox="`0 0 ${media.dim_width} ${media.dim_height}`"
      class="absolute inset-0 text-yellow-300 text-opacity-50 pointer-events-none fill-current"
    >
      <rect v-for="b in blocks" :key="b.id" :x="b.min_x" :y="b.min_y" :width="b.max_x - b.min_x" :height="b.max_y - b.min_y" ry="8" />
    </svg>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { urlMedia, MediaType } from '~/request'
import { VideoCameraIcon } from '@heroicons/vue/outline'
import useStore from '~/composables/useStore'

const props = defineProps<{
  hash?: string
  thumb?: boolean
}>()

const emit = defineEmits<{ (e: 'loaded'): void }>()
const loaded = () => emit('loaded')

const store = useStore()

const media = computed(() => store.getMediaByHash(props.hash || ''))
const blocks = computed(() => store.getHighlightedBlocksByHash(props.hash || ''))
const url = computed(
  () =>
    props.thumb
      ? `${urlMedia}/${media.value?.hash}/thumb` // ~200px thumb
      : `${urlMedia}/${media.value?.hash}/raw`, // full image or video
)

const isVideo = computed(() => media.value?.type === MediaType.Video)
</script>
