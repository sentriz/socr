<template>
  <div
    v-if="screenshot"
    class="fixed h-full top-0 right-0 w-9/12 border-l-4 p-6 bg-white"
  >
    <div class="mx-auto">
      <div class="bg-black shadow font-mono text-sm">
        <ScreenshotHighlight class="mx-auto" :screenshot="screenshot" />
      </div>
      <hr />
      <div class="bg-gray-300 padded shadow font-mono text-sm">
        <p v-for="(line, i) in text" :key="i">
          {{ line }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup="props">
export { default as ScreenshotHighlight } from "./ScreenshotHighlight.vue";
export default {
  props: {
    id: String,
    results: Array,
  },
};

import { computed } from "vue";
import { fields } from "../api/";

// TODO: not pass all results to this component
// perhaps use vuex
export const screenshot = computed(() =>
  props.results.find((result) => result.id === props.id)
);
export const text = computed(() => screenshot.value.fields[fields.BLOCKS_TEXT]);
</script>
