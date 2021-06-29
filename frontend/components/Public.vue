<template>
  <div class="bg-gray-50 flex items-center min-h-screen">
    <div class="container p-6 mx-auto space-y-6">
      <!-- block image or loading -->
      <media-background v-show="mediaHave" :hash="media?.hash" class="box p-3">
        <media-highlight rich :hash="media?.hash" class="mx-auto" @loaded="mediaLoaded" />
      </media-background>
      <loading-spinner v-show="!mediaHave" class="bg-gray-100" text="processing image" />
      <!-- block text or loading -->
      <template v-if="!isVideo">
        <div v-show="blocks.length" class="box padded font-mono text-sm bg-white">
          <p v-for="(block, i) in blocks" :key="i">
            {{ block.body }}
          </p>
        </div>
        <loading-spinner v-show="!blocks.length" class="bg-gray-100" text="processing text" />
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import MediaBackground from './MediaBackground.vue'
import MediaHighlight from './MediaHighlight.vue'
import LoadingSpinner from './LoadingSpinner.vue'

import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { newSocket, MediaType } from '../api'
import useStore from '../composables/useStore'

const store = useStore()
const route = useRoute()
const hash = (route.params.hash as string) || ''

const mediaHave = ref(false)
const mediaLoaded = () => (mediaHave.value = true)

const media = computed(() => store.getMediaByHash(hash))
const blocks = computed(() => store.getBlocksByHash(hash))
const isVideo = computed(() => media.value?.type === MediaType.Video)

const requestMedia = async () => {
  await store.loadMedia(hash)
}

// fetch image on mount
onMounted(requestMedia)

// fetch image on socket message
const socket = newSocket({ want_media_hash: hash })
socket.onmessage = requestMedia
</script>
