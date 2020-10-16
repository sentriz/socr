<template>
  <Transition
    enterActiveClass="transform transition ease-in-out duration-200"
    enterToClass="translate-x-0"
    enterFromClass="translate-x-full"
    leaveActiveClass="transform transition ease-in-out duration-200"
    leaveToClass="translate-x-full"
    leaveFromClass="translate-x-0"
  >
    <div
      v-if="screenshot"
      class="z-20 fixed inset-y-0 right-0 w-9/12 p-6 bg-white overflow-y-auto"
    >
      <div class="mx-auto space-y-6">
        <div class="bg-black shadow">
          <ScreenshotHighlight class="mx-auto" :id="screenshot.id" />
        </div>
        <div class="bg-gray-300 shadow font-mono text-sm padded">
          <p v-for="(line, i) in text" :key="i">
            {{ line }}
          </p>
        </div>
      </div>
    </div>
  </Transition>
  <Transition
    enterActiveClass="ease-in-out duration-500"
    enterToClass="opacity-100"
    enterFromClass="opacity-0"
    leaveActiveClass="ease-in-out duration-500"
    leaveToClass="opacity-0"
    leaveFromClass="opacity-100"
  >
    <div
      v-if="screenshot"
      class="z-10 fixed inset-0 bg-gray-700 bg-opacity-50 transition-opacity pointer-events-none"
    />
  </Transition>
</template>

<script setup="props">
import ScreenshotHighlight from "./ScreenshotHighlight.vue";
export default {
  components: { ScreenshotHighlight },
  props: { id: String },
};

import { inject, computed, watch } from "vue";
import { fields } from "../api/";
import { useStore } from "../store/";

export const store = useStore();

// load the screenshot from the network if we can't find it in the store
// (can happen on page reload if we've click an image on the eg. 5th page)
watch(
  () => props.id,
  (id) => {
    if (id && !store.screenshotByID(id)) {
      store.screenshotsLoadID(id);
    }
  },
  { immediate: true },
);

export const screenshot = computed(() => store.screenshotByID(props.id));
export const text = computed(() => screenshot.value.fields[fields.BLOCKS_TEXT]);
</script>
