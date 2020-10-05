<template>
  <div class="fixed inset-0 flex pointer-events-none">
    <div class="w-3/12 bg-gradient-to-l from-black to-transparent opacity-50" />
    <div class="w-9/12 border-l-4 p-6 bg-white overflow-y-auto">
      <div class="mx-auto">
        <div class="my-2 bg-black shadow">
          <ScreenshotHighlight class="mx-auto" :id="screenshot.id" />
        </div>
        <hr />
        <div class="my-2 bg-gray-300 shadow font-mono text-sm padded">
          <p v-for="(line, i) in text" :key="i">
            {{ line }}
          </p>
        </div>
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
