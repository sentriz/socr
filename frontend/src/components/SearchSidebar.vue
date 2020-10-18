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
      class="z-20 fixed inset-y-0 right-0 w-9/12 p-6 bg-gray-100 overflow-y-auto"
    >
      <div class="space-y-6">
        <div class="text-right space-x-2">
          <span>
            created
            <span
              class="badge bg-pink-200 text-pink-900"
              :title="screenshot.fields.timestamp"
            >
              {{ relativeDateStr(screenshot.fields.timestamp) }}
            </span>
          </span>
          <span v-if="tags?.length">
            tags
            <span class="space-x-2">
              <span
                class="badge bg-blue-200 text-blue-900"
                v-for="(tag, i) in tags"
                :id="i"
              >
                {{ tag }}
              </span>
            </span>
          </span>
        </div>
        <div class="box bg-white">
          <ScreenshotHighlight class="mx-auto max-w-full" :id="screenshot.id" />
        </div>
        <div class="box bg-gray-200 padded font-mono text-sm">
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
      class="z-10 fixed inset-0 bg-gray-700 bg-opacity-75 transition-opacity"
    >
      <div class="w-3/12 p-6 flex justify-end text-white text-2xl pointer-events-none">
        <router-link :to="{ name: 'search' }" class="pointer-events-auto">
          <i class="fas fa-times-circle"></i>
        </router-link>
      </div>
    </div>
  </Transition>
</template>

<script setup="props">
import ScreenshotHighlight from "./ScreenshotHighlight.vue";
export default {
  components: { ScreenshotHighlight },
  props: { id: String },
};

import { inject, computed, watch } from "vue";
import relativeDate from "relative-date";
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

export const relativeDateStr = (stamp) => relativeDate(new Date(stamp));

export const screenshot = computed(() => store.screenshotByID(props.id));
export const text = computed(() => screenshot.value.fields[fields.BLOCKS_TEXT]);
export const tags = computed(() => {
  const tags = screenshot.value.fields[fields.TAGS];
  if (tags) return Array.isArray(tags) ? tags : [tags];
});
</script>
