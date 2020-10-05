<template>
  <div
    v-if="screenshot"
    class="fixed h-full top-0 right-0 w-9/12 border-l-4 p-6 bg-white overflow-y-auto"
  >
    <div class="mx-auto">
      <div class="bg-black shadow font-mono text-sm">
        <ScreenshotHighlight class="mx-auto" :id="screenshot.id" />
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
  },
};

import { inject, computed, watch } from "vue";
import { fields } from "../api/";

export const store = inject("store");
export const screenshot = computed(() => store.screenshots[props.id]);
export const text = computed(() => screenshot.value.fields[fields.BLOCKS_TEXT]);
</script>
