<template>
  <div class="flex min-h-screen flex-col justify-center bg-white">
    <div class="container mx-auto p-5">
      <div class="flex items-center justify-between">
        <h1 class="text-gray-700">shared media</h1>
        <badge-group label="created on" class="hidden text-gray-500 md:inline-flex" v-if="media && timestamp">
          <badge class="bg-pink-200 text-pink-900" :title="media.timestamp">{{ timestamp }}</badge>
        </badge-group>
      </div>
    </div>
    <media-preview :hash="hash" class="max-h-[750px] py-2" />
    <div v-if="!isVideo" class="container mx-auto p-5">
      <media-lines :hash="hash" />
    </div>
  </div>
</template>

<script setup lang="ts">
import MediaPreview from './MediaPreview.vue'
import MediaLines from './MediaLines.vue'
import BadgeGroup from './BadgeGroup.vue'
import Badge from './Badge.vue'
import { MediaType } from '~/request'
import { computed, onMounted } from 'vue'
import { newSocket } from '~/request'
import { useRoute } from 'vue-router'
import useStore from '~/composables/useStore'

const store = useStore()
const route = useRoute()
const hash = (route.params.hash as string) || ''

const requestMedia = async () => {
  await store.loadMedia(hash)
}

const media = computed(() => store.getMediaByHash(hash))
const timestamp = computed(() => {
  if (!media.value) return
  const date = new Date(media.value?.timestamp)
  return `${pad(date.getFullYear())}.${pad(date.getMonth())}.${pad(date.getDay())}`
})

const pad = (n: number): string => String(n).padStart(2, '0')

onMounted(requestMedia)

const isVideo = computed(() => media.value?.type === MediaType.Video)

const socket = newSocket({ want_media_hash: hash })
socket.onmessage = requestMedia
</script>
