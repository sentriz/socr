<template>
  <transition-fade>
    <div v-if="media" class="fixed inset-0 z-10 transition-opacity bg-gray-700 bg-opacity-75" />
  </transition-fade>
  <transition-slide>
    <div v-if="media" ref="content" class="fixed inset-y-0 right-0 z-20 w-full max-w-lg">
      <div class="bg-gray-50 h-full p-6 space-y-6 overflow-y-auto">
        <search-sidebar-header :hash="media.hash" />
        <media-preview :hash="media.hash" />
      </div>
    </div>
  </transition-slide>
</template>

<script setup lang="ts">
import TransitionFade from './TransitionFade.vue'
import TransitionSlide from './TransitionSlideX.vue'
import SearchSidebarHeader from './SearchSidebarHeader.vue'
import MediaPreview from './MediaPreview.vue'

import { computed, ref, watch } from 'vue'
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

const content = ref<HTMLElement>()
onClickOutside(content, () => router.push({ name: 'search' }))
</script>
