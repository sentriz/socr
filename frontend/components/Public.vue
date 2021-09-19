<template>
  <div class="bg-gray-50 flex items-center min-h-screen">
    <div class="container p-6 mx-auto space-y-6">
      <media-preview :hash="hash" />
    </div>
  </div>
</template>

<script setup lang="ts">
import MediaPreview from './MediaPreview.vue'
import { onMounted } from 'vue'
import { newSocket } from '../api'
import { useRoute } from 'vue-router'
import useStore from '../composables/useStore'

const store = useStore()
const route = useRoute()
const hash = (route.params.hash as string) || ''

const requestMedia = async () => {
  await store.loadMedia(hash)
}

onMounted(requestMedia)

const socket = newSocket({ want_media_hash: hash })
socket.onmessage = requestMedia
</script>
