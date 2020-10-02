<!-- <ScreenshotHighlight v-if="screenshot" :screenshot="screenshot" /> -->
<template>
  <div class="bg-gray-200 min-h-screen">
    <div v-show="imageHave" class="container mx-auto p-8 space-y-4">
      <div class="box bg-black">
        <img class="mx-auto" :src="imageSrc" @load="imageLoaded" />
      </div>
      <div class="box bg-gray-300 padded font-mono text-sm">
        <p v-show="text.length == 0" class="text-gray-600 py-2">
          <i class="animate-spin fas fa-circle-notch"></i> processing screenshot...
        </p>
        <p v-for="(line, i) in text" :key="i">
          {{ line }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
export { default as ScreenshotHighlight } from "./ScreenshotHighlight.vue";
export default {
  props: {},
};

import { ref, reactive, watch, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { reqScreenshot, fields, urlScreenshot, newSocket } from "../api";

const router = useRouter();
const route = useRoute();
const screenshotID = route.params.id;

export const screenshot = ref(null);
export const text = ref([]);
const requestScreenshot = async () => {
  const resp = await reqScreenshot(screenshotID);
  if (resp.hits.length == 0) return;

  const hit = resp.hits[0];
  screenshot.value = hit;
  text.value = hit.fields[fields.BLOCKS_TEXT];
};

// suspend showing anything until we have an image
export const imageSrc = `${urlScreenshot}/${screenshotID}/raw`;
export const imageHave = ref(false);
export const imageLoaded = () => {
  imageHave.value = true;
};

// fetch image on mount
onMounted(requestScreenshot);

// fetch image on socket message
const socket = newSocket({ want_screenshot_id: screenshotID });
socket.onmessage = requestScreenshot;
</script>
