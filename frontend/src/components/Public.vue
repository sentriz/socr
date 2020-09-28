<template>
  <div class="container">
    <ScreenshotHighlight
      v-if="screenshot"
      :screenshot="screenshot"
      class="photo border border-gray-300 rounded-lg"
    />
  </div>
</template>

<script setup>
export { default as ScreenshotHighlight } from "./ScreenshotHighlight.vue";

import { ref, reactive, watch, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { reqSearch, fields } from "../api";

const route = useRoute();
const screenshotID = route.params.id;

export const screenshot = ref(null);
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
  }
});
</script>
