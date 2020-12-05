<template>
  <Transition
    enter-active-class="transform transition ease-in-out duration-200"
    enter-to-class="translate-x-0"
    enter-from-class="translate-x-full"
    leave-active-class="transform transition ease-in-out duration-200"
    leave-to-class="translate-x-full"
    leave-from-class="translate-x-0"
  >
    <div v-if="screenshot" class="z-20 fixed inset-y-0 right-0 max-w-lg w-full p-6 bg-gray-200 overflow-y-auto">
      <!-- sidebar main column -->
      <div class="space-y-6">
        <!-- sidebar header -->
        <div class="flex leading-normal">
          <router-link :to="{ name: 'search' }" class="flex-grow text-xl leading-none">
            <i class="text-gray-900 fas fa-times-circle"></i>
          </router-link>
          <!-- sidebar header badges -->
          <div class="flex flex-col md:flex-row gap-3 justify-end items-end">
            <BadgeLabel label="created">
              <Badge class="badge bg-pink-200 text-pink-900" :title="timestamp">
                {{ relativeDateStr(timestamp) }}
              </Badge>
            </BadgeLabel>
            <BadgeLabel v-if="tags?.length" label="tags">
              <Badge v-for="(tag, i) in tags" :key="i" class="badge bg-blue-200 text-blue-900">
                {{ tag }}
              </Badge>
            </BadgeLabel>
            <Badge class="bg-indigo-200 text-indigo-900" icon="fas fa-external-link-alt">
              <a :href="screenshotRaw" target="_blank">raw</a>
            </Badge>
            <Badge class="bg-green-200 text-green-900" icon="fas fa-external-link-alt">
              <router-link :to="{ name: 'public', params: { id: screenshot.id } }"> public </router-link>
            </Badge>
          </div>
        </div>
        <!-- sidebar box -->
        <ScreenshotBackground :id="screenshot.id" class="box p-3">
          <ScreenshotHighlight :id="screenshot.id" class="mx-auto" />
        </ScreenshotBackground>
        <!-- sidebar box -->
        <div v-if="text.length" class="box bg-gray-100 padded font-mono text-sm">
          <p v-for="(line, i) in text" :key="i">
            {{ line }}
          </p>
        </div>
      </div>
    </div>
  </Transition>
  <Transition
    enter-active-class="ease-in-out duration-500"
    enter-to-class="opacity-100"
    enter-from-class="opacity-0"
    leave-active-class="ease-in-out duration-500"
    leave-to-class="opacity-0"
    leave-from-class="opacity-100"
  >
    <div
      v-if="screenshot"
      class="z-10 fixed inset-0 bg-gray-700 bg-opacity-75 transition-opacity pointer-events-none"
    />
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
