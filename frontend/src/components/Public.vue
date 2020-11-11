<template>
  <div class="bg-gray-200 min-h-screen">
    <div v-show="imageHave" class="container mx-auto p-6 space-y-6 flex flex-col">
      <ScreenshotBackground :id="screenshot?.id" class="box p-3">
        <img class="mx-auto" :src="imageSrc" @load="imageLoaded" />
      </ScreenshotBackground>
      <div class="box bg-gray-100 padded font-mono text-sm">
        <p v-show="text.length == 0" class="text-gray-600 py-2">
          <i class="animate-spin fas fa-circle-notch"></i>
          processing screenshot...
        </p>
        <p v-for="(line, i) in text" :key="i">
          {{ line }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup="props">
import ScreenshotBackground from "./ScreenshotBackground.vue";
export default {
  components: { ScreenshotBackground },
  props: {},
};

import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { fields as apifields, urlScreenshot, newSocket } from "../api";
import { useStore } from "../store";

const store = useStore();
const route = useRoute();
const screenshotID = route.params.id;

export const screenshot = ref(null);
const requestScreenshot = async () => {
  const resp = await store.screenshotsLoadID(screenshotID);
  if (resp.hits.length == 0) return;

  screenshot.value = store.screenshotByID(screenshotID);
};

export const fields = computed(() => screenshot.value?.fields);
export const text = computed(() => fields.value?.[apifields.BLOCKS_TEXT] || []);
export const timestamp = computed(() => fields.value?.[apifields.TIMESTAMP] || "");

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
