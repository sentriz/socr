<template>
  <input
    class="inp w-full"
    type="text"
    placeholder="enter screenshot text query"
    v-model="query"
  />
  <p class="my-3 text-gray-500 text-right">
    {{ reqTotalHits }} results found in {{ reqTookMs }}ms
  </p>
  <div class="col-resp gap-x-3 space-y-3">
    <div v-for="screenshot in store.screenshots" :key="screenshot.id">
      <router-link :to="{ name: 'result', params: { id: screenshot.id } }">
        <ScreenshotHighlight :id="screenshot.id" class="border border-gray-300 rounded" />
      </router-link>
    </div>
  </div>
  <teleport to="body">
    <router-view v-slot="{ Component, route }">
      <Transition
        enterFromClass="translate-x-full"
        enterActiveClass="transform transition ease-in-out duration-200"
        enterToClass="translate-x-0"
        leaveFromClass="translate-x-0"
        leaveActiveClass="transform transition ease-in-out duration-200"
        leaveToClass="translate-x-full"
      >
        <component :is="Component" v-if="reqTotalHits" v-bind="route.params"></component>
      </Transition>
    </router-view>
  </teleport>
</template>

<script setup>
export { default as ScreenshotHighlight } from "./ScreenshotHighlight.vue";
export default {
  props: {},
};

import { inject, ref, reactive, watch, computed } from "vue";
import throttle from "lodash.debounce";
import { reqSearch, fields } from "../api";

export const query = ref("");
watch(query, (query, _) => {
  if (query) fetchScreenshots();
});

const searchParams = {
  size: 40,
  fields: [
    fields.BLOCKS_TEXT,
    fields.BLOCKS_POSITION,
    fields.SIZE_HEIGHT,
    fields.SIZE_WIDTH,
  ],
  highlight: {
    fields: [fields.BLOCKS_TEXT],
  },
};

export const store = inject("store");
export const reqTotalHits = ref(0);
export const reqTookMs = ref(0);
export const fetchScreenshots = throttle(async () => {
  const resp = await reqSearch({
    ...searchParams,
    query: { term: query.value },
  });

  reqTotalHits.value = resp.total_hits;
  reqTookMs.value = Math.round((resp.took / 100000) * 100) / 100;

  store.screenshots = {};
  for (const hit of resp.hits) {
    store.screenshots[hit.id] = hit;
  }
}, 200);
</script>
