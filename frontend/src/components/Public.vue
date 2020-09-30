<!-- <ScreenshotHighlight v-if="screenshot" :screenshot="screenshot" /> -->
<template>
  <div class="bg-gray-200 h-full">
    <div class="container mx-auto p-8 space-y-4">
      <img
        class="border-2 border-solid border-white shadow mx-auto"
        :src="url"
      />
      <div
        class="border-2 border-solid border-white shadow bg-gray-300 padded font-mono text-sm"
      >
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
import { useRoute } from "vue-router";
import { reqImage, fields, urlImage, newSocket } from "../api";

const route = useRoute();
const screenshotID = route.params.id;

export const screenshot = ref(null);
export const text = ref([]);
export const url = ref("");

const setStateFromScreenshot = (scrot, x) => {
  console.log(x, "/", scrot);
  screenshot.value = scrot;
  text.value = scrot.fields[fields.BLOCKS_TEXT];
  url.value = `${urlImage}/${scrot.id}/raw`;
};

// screenshot data from xhr
onMounted(async () => {
  const resp = await reqImage(screenshotID);
  if (resp.hits.length > 0) setStateFromScreenshot(resp.hits[0], "a");
});

// screenshot data from socket
const socket = newSocket({ want_screenshot_id: screenshotID });
socket.onmessage = (e) => {
  setStateFromScreenshot(JSON.parse(e.data), "b");
};
</script>
