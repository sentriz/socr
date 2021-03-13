<template>
  <div class="bg-gray-100 min-h-screen">
    <div class="container mx-auto p-6 flex flex-col gap-6">
      <!-- block image or loading -->
      <ScreenshotBackground v-show="imageHave" :hash="screenshot?.hash || ''" class="box p-3">
        <img class="mx-auto" :src="image" @load="imageLoaded" />
      </ScreenshotBackground>
      <LoadingSpinner v-show="!imageHave" class="bg-gray-100" text="processing image" />
      <!-- block text or loading -->
      <div v-show="blocks.length" class="box bg-white padded font-mono text-sm">
        <p v-for="(block, i) in blocks" :key="i">
          {{ block.body }}
        </p>
      </div>
      <LoadingSpinner v-show="!blocks.length" class="bg-gray-100" text="processing text" />
    </div>
  </div>
</template>

<script setup lang="ts">
import ScreenshotBackground from './ScreenshotBackground.vue'
import LoadingSpinner from './LoadingSpinner.vue'

import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { urlScreenshot, newSocket } from '../api'
import useStore from '../composables/useStore'

const store = useStore()
const route = useRoute()
const hash = (route.params.hash as string) || ''

const image = ref("")
const imageHave = ref(false)
const imageLoaded = (_: Event) => {
  imageHave.value = true
}

const screenshot = computed(() => store.getScreenshotByHash(hash))
const blocks = computed(() => store.getBlocksByHash(hash))

const requestScreenshot = async () => {
  const now = new Date()
  image.value = `${urlScreenshot}/${hash}/raw?t=${now.valueOf()}`

  await store.loadScreenshot(hash)
}

// fetch image on mount
onMounted(requestScreenshot)

// fetch image on socket message
const socket = newSocket({ want_screenshot_hash: hash })
socket.onmessage = requestScreenshot
</script>
