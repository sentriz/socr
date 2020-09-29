<!-- <ScreenshotHighlight v-if="screenshot" :screenshot="screenshot" /> -->
<template>
  <div class="bg-gray-200 h-full">
    <div class="container mx-auto p-8 space-y-4">
      <img
        class="border-2 border-solid border-white shadow mx-auto"
        src="/api/image/1MzJAZZA7d5blOeOKV6jQzyeM4qECfjw"
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
import { reqSearch, fields } from "../api";

const route = useRoute();
const screenshotID = route.params.id;

export const screenshot = ref(null);
export const text = ref([]);
onMounted(async () => {
  const resp = await reqSearch({
    fields: [
      fields.BLOCKS_TEXT,
      fields.BLOCKS_POSITION,
      fields.SIZE_HEIGHT,
      fields.SIZE_WIDTH,
    ],
    highlight: {
      fields: [fields.BLOCKS_TEXT],
    },
    query: {
      ids: [screenshotID],
    },
  });

  if (resp.hits.length > 0) {
    screenshot.value = resp.hits[0];
    text.value = resp.hits[0].fields[fields.BLOCKS_TEXT];
  }
});
</script>
