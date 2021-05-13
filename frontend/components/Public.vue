<template>
  <div class="bg-gray-50 min-h-screen">
    <div class="container flex flex-col gap-6 p-6 mx-auto">
      <!-- block image or loading -->
      <media-background v-show="imageHave" :hash="media?.hash || ''" class="box p-3">
        <img class="mx-auto" :src="image" @load="imageLoaded" />
      </media-background>
      <loading-spinner v-show="!imageHave" class="bg-gray-100" text="processing image" />
      <!-- block text or loading -->
      <div v-show="blocks.length" class="box padded font-mono text-sm bg-white">
        <p v-for="(block, i) in blocks" :key="i">
          {{ block.body }}
        </p>
      </div>
      <loading-spinner v-show="!blocks.length" class="bg-gray-100" text="processing text" />
    </div>
  </div>
</template>

<script setup lang="ts">
import MediaBackground from './MediaBackground.vue'
import LoadingSpinner from './LoadingSpinner.vue'

import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { urlMedia, newSocket } from '../api'
import useStore from '../composables/useStore'

const store = useStore()
const route = useRoute()
const hash = (route.params.hash as string) || ''

const image = ref('')
const imageHave = ref(false)
const imageLoaded = (_: Event) => {
  imageHave.value = true
}

const media = computed(() => store.getMediaByHash(hash))
const blocks = computed(() => store.getBlocksByHash(hash))

const requestMedia = async () => {
  const now = new Date()
  image.value = `${urlMedia}/${hash}/raw?t=${now.valueOf()}`

  await store.loadMedia(hash)
}

// fetch image on mount
onMounted(requestMedia)

// fetch image on socket message
const socket = newSocket({ want_media_hash: hash })
socket.onmessage = requestMedia
</script>
