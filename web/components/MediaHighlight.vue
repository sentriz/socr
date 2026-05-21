<template>
  <div class="relative" v-bind="$attrs">
    <div v-if="thumb && isVideo" class="absolute right-2 top-2 rounded-md bg-black/80 p-2">
      <video-camera-icon class="h-5 text-white" />
    </div>
    <video :key="url" v-if="!thumb && isVideo && media" :controls="true" @loadstart="loaded" class="block w-full max-h-[inherit]">
      <source :src="url" :type="media.mime" />
    </video>
    <img v-else :src="url" @load="loaded" class="block w-full max-h-[inherit] object-cover" />

    <svg v-if="media && blocks.length" :viewBox="`0 0 ${media.dim_width} ${media.dim_height}`" class="pointer-events-none absolute inset-0 fill-amber-300/50">
      <rect v-for="b in blocks" :key="b.id" :x="b.min_x" :y="b.min_y" :width="b.max_x - b.min_x" :height="b.max_y - b.min_y" ry="4" />
    </svg>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { urlMedia, MediaType } from '~/request'
import { VideoCameraIcon } from '@heroicons/vue/24/outline'
import useStore from '~/composables/useStore'

defineOptions({
  inheritAttrs: false,
})

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
