<template>
  <div class="flex min-h-screen items-center bg-white">
    <div class="container mx-auto space-y-6 p-6">
      <div v-if="media" class="flex items-center justify-between">
        <h1 class="text-gray-500">shared media</h1>
        <badge-group label="uploaded" class="hidden text-gray-700 md:inline-flex" v-if="media && timestamp">
          <badge class="bg-pink-200 text-pink-900" :title="media.timestamp">{{ timestamp }}</badge>
        </badge-group>
      </div>
      <media-preview :hash="hash" />
    </div>
  </div>
</template>

<script setup lang="ts">
import MediaPreview from './MediaPreview.vue'
import BadgeGroup from './BadgeGroup.vue'
import Badge from './Badge.vue'
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
  return new Date(media.value?.timestamp).toLocaleString()
})

onMounted(requestMedia)

const socket = newSocket({ want_media_hash: hash })
socket.onmessage = requestMedia
</script>
