<template>
  <Transition
    enter-active-class="transform transition ease-in-out duration-200"
    enter-to-class="translate-x-0"
    enter-from-class="translate-x-full"
    leave-active-class="transform transition ease-in-out duration-200"
    leave-to-class="translate-x-full"
    leave-from-class="translate-x-0"
  >
  <SearchSidebarMain :id="screenshot?.id || ''" />
  </Transition>
  <Transition
    enter-active-class="ease-in-out duration-500"
    enter-to-class="opacity-100"
    enter-from-class="opacity-0"
    leave-active-class="ease-in-out duration-500"
    leave-to-class="opacity-0"
    leave-from-class="opacity-100"
  >
    <SearchSidebarBackground :id="screenshot?.id || ''" />
  </Transition>
</template>

<script setup lang="ts">
import ScreenshotHighlight from "./ScreenshotHighlight.vue";
import ScreenshotBackground from "./ScreenshotBackground.vue";
import SearchSidebarMain from "./SearchSidebarMain.vue"
import SearchSidebarBackground from "./SearchSidebarBackground.vue"
import BadgeLabel from "./BadgeLabel.vue";
import Badge from "./Badge.vue";

import { computed, defineProps, watch } from "vue";
import relativeDate from "relative-date";
import { urlScreenshot, Field } from "../api/";
import type { Store } from "../store"
import useStore from "../composables/useStore";

const props = defineProps<{
  id: string | undefined,
}>();

const store = useStore() || {} as Store;

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

const relativeDateStr = (stamp: string) => relativeDate(new Date(stamp));

const screenshotRaw = computed(() => `${urlScreenshot}/${props.id}/raw`);
const screenshot = computed(() => store.screenshotByID(props.id || ""));
const text = computed(() => toArray(screenshot.value?.fields[Field.BLOCKS_TEXT] || []));
const timestamp = computed(() => `${screenshot.value?.fields[Field.TIMESTAMP]}`);
const tags = computed(() => {
  const tags = screenshot.value.fields[Field.TAGS];
  if (tags) return Array.isArray(tags) ? tags : [tags];
});

const toArray = <T>(value: T | T[]) => (Array.isArray(value) ? value : [value]);
</script>
