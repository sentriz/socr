<template>
  <div class="bg-gray-200 min-h-screen">
    <div v-show="imageHave" class="container mx-auto p-6 flex flex-col gap-6">
      <ScreenshotBackground :hash="screenshot?.hash || ''" class="box p-3">
        <img class="mx-auto" :src="imageSrc" @load="imageLoaded" />
      </ScreenshotBackground>
      <div class="box bg-gray-100 padded font-mono text-sm">
        <p v-show="blocks.length == 0" class="text-gray-600 py-2">
          <i class="animate-spin fas fa-circle-notch"></i>
          processing screenshot...
        </p>
        <p v-for="(block, i) in blocks" :key="i">
          {{ block.body }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ScreenshotBackground from "./ScreenshotBackground.vue";

import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { urlScreenshot, newSocket, isError } from "../api";
import type { Screenshot } from "../api";
import useStore from "../composables/useStore";

const store = useStore();
const route = useRoute();
const hash = route.params.hash as string || "";

const screenshot = ref<Screenshot>();
const requestScreenshot = async () => {
  const resp = await store.loadScreenshot(hash);
  if (isError(resp)) return

  screenshot.value = store.getScreenshotByHash(hash);
};

const blocks = computed(() => store.getBlocksByHash(hash));

// suspend showing anything until we have an image
const imageSrc = `${urlScreenshot}/${hash}/raw`;
const imageHave = ref(false);
const imageLoaded = () => {
  imageHave.value = true;
};

// fetch image on mount
onMounted(requestScreenshot);

// fetch image on socket message
const socket = newSocket({ want_screenshot_hash: hash });
socket.onmessage = requestScreenshot;
</script>
