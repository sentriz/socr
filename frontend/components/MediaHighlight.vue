<template>
  <div class="w-fit relative">
    <video v-if="rich && media?.type === MediaType.Video" :controls="true" @loadstart="loaded">
      <source :src="url" :type="media.mime" />
    </video>
    <img v-else :src="url" @load="loaded" />
    <svg
      v-if="media && blocks.length"
      :viewBox="`0 0 ${media.dim_width} ${media.dim_height}`"
      class="absolute inset-0 text-yellow-500 text-opacity-50 pointer-events-none fill-current"
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
    <div v-if="!rich && media?.type === MediaType.Video" class="top-2 right-2 bg-black/80 absolute p-2 rounded-md">
      <video-camera-icon class="h-5 text-white" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { urlMedia, MediaType } from '../api'
import { VideoCameraIcon } from '@heroicons/vue/outline'
import useStore from '../composables/useStore'

const props = defineProps<{
  hash?: string
  rich?: boolean
}>()

const emit = defineEmits<{ (e: 'loaded'): void }>()
const loaded = () => emit('loaded')

const store = useStore()

const media = computed(() => store.getMediaByHash(props.hash || ''))
const blocks = computed(() => store.getHighlightedBlocksByHash(props.hash || ''))
const url = computed(
  () =>
    props.rich
      ? `${urlMedia}/${media.value?.hash}/raw` // full image or video
      : `${urlMedia}/${media.value?.hash}/thumb`, // ~200px thumb
)
</script>
