<template>
  <transition-fade>
    <div v-if="media" class="pointer-events-auto absolute inset-0 bg-gray-700 bg-opacity-75 transition-opacity" />
  </transition-fade>
  <transition-slide>
    <div v-if="media" ref="content" class="overflow-y-thin pointer-events-auto absolute inset-y-0 right-0 w-full max-w-lg space-y-6 bg-white p-6">
      <search-sidebar-header :hash="media.hash" />
      <media-preview :hash="media.hash" />
      <media-lines v-if="!isVideo" :hash="media.hash" />
    </div>
  </transition-slide>
</template>

<script setup lang="ts">
import TransitionFade from './TransitionFade.vue'
import TransitionSlide from './TransitionSlideX.vue'
import SearchSidebarHeader from './SearchSidebarHeader.vue'
import MediaLines from './MediaLines.vue'
import MediaPreview from './MediaPreview.vue'
import { MediaType } from '~/request'

import { computed, ref, watch } from 'vue'
import useStore from '~/composables/useStore'
import { onClickOutside } from '@vueuse/core'
import { useRouter } from 'vue-router'
import { routes } from '~/router'

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
const isVideo = computed(() => media.value?.type === MediaType.Video)

const content = ref<HTMLElement>()
onClickOutside(content, () => router.push({ name: routes.SEARCH }))
</script>
