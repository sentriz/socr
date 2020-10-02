<!-- <ScreenshotHighlight v-if="screenshot" :screenshot="screenshot" /> -->
<template>
  <div class="bg-gray-200 min-h-screen">
    <div v-show="screenshot" class="container mx-auto p-8 space-y-4">
      <div class="border-2 border-solid border-white bg-black shadow min-h-2">
        <img class="mx-auto" :src="url" />
      </div>
      <div
        class="border-2 border-solid border-white shadow bg-gray-300 padded font-mono text-sm"
      >
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

export const url = `${urlScreenshot}/${screenshotID}/raw`;
export const screenshot = ref(null);
export const text = ref([]);

const requestScreenshot = async () => {
  const resp = await reqScreenshot(screenshotID);
  if (resp.hits.length == 0) {
    router.replace({ name: "not_found" });
    return;
  }

  const hit = resp.hits[0];
  screenshot.value = hit;
  text.value = hit.fields[fields.BLOCKS_TEXT];
};

// fetch image on mount
onMounted(requestScreenshot);

// fetch image on socket message
const socket = newSocket({ want_screenshot_id: screenshotID });
socket.onmessage = requestScreenshot;
</script>
