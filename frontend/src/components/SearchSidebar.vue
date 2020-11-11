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
      class="z-20 fixed inset-y-0 right-0 w-9/12 p-6 bg-gray-200 overflow-y-auto"
    >
      <div class="space-y-6">
        <div
          class="flex flex-col md:flex-row space-y-2 md:space-x-4 justify-end items-end"
        >
          <BadgeLabel label="created">
            <Badge class="badge bg-pink-200 text-pink-900" :title="timestamp">
              {{ relativeDateStr(timestamp) }}
            </Badge>
          </BadgeLabel>
          <BadgeLabel v-if="tags?.length" label="tags">
            <Badge v-for="tag in tags" class="badge bg-blue-200 text-blue-900">
              {{ tag }}
            </Badge>
          </BadgeLabel>
          <Badge class="bg-indigo-200 text-indigo-900" icon="fas fa-external-link-alt">
            <a :href="screenshotRaw" target="_blank">raw</a>
          </Badge>
          <Badge class="bg-green-200 text-green-900" icon="fas fa-external-link-alt">
            <router-link :to="{ name: 'public', params: { id: screenshot.id } }">
              public
            </router-link>
          </Badge>
        </div>
        <ScreenshotBackground :id="screenshot.id" class="box p-3">
          <ScreenshotHighlight :id="screenshot.id" class="mx-auto" />
        </ScreenshotBackground>
        <div v-if="text.length" class="box bg-gray-100 padded font-mono text-sm">
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
      <div class="w-3/12 p-6 flex justify-end text-white text-xl pointer-events-none">
        <div>
          <router-link :to="{ name: 'search' }" class="pointer-events-auto">
            <i class="fas fa-times-circle"></i>
          </router-link>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup="props">
import ScreenshotHighlight from "./ScreenshotHighlight.vue";
import ScreenshotBackground from "./ScreenshotBackground.vue";
import BadgeLabel from "./BadgeLabel.vue";
import Badge from "./Badge.vue";
export default {
  components: {
    ScreenshotHighlight,
    ScreenshotBackground,
    BadgeLabel,
    Badge,
  },
  props: { id: String },
};

import { computed, watch } from "vue";
import relativeDate from "relative-date";
import { urlScreenshot, fields } from "../api/";
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

export const screenshotRaw = computed(() => `${urlScreenshot}/${props.id}/raw`);
export const screenshot = computed(() => store.screenshotByID(props.id));
export const text = computed(() => screenshot.value.fields[fields.BLOCKS_TEXT]);
export const timestamp = computed(() => screenshot.value.fields[fields.TIMESTAMP]);
export const tags = computed(() => {
  const tags = screenshot.value.fields[fields.TAGS];
  if (tags) return Array.isArray(tags) ? tags : [tags];
});
</script>
