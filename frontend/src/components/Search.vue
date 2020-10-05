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
  <hr class="my-0" />
  <div id="photos">
    <router-link
      v-for="screenshot in store.screenshots"
      :key="screenshot.id"
      :to="{ name: 'result', params: { id: screenshot.id } }"
    >
      <ScreenshotHighlight
        :id="screenshot.id"
        class="photo border border-gray-300 rounded-lg"
      />
    </router-link>
  </div>
  <teleport to="body">
    <router-view v-slot="{ Component, route }">
      <transition name="sidebar-slide">
        <component :is="Component" v-if="reqTotalHits" v-bind="route.params"></component>
      </transition>
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

<style scoped>
#photos {
  line-height: 0;
  column-count: 4;
  column-gap: 5px;
}

#photos .photo {
  width: 100%;
  height: auto;
  margin: 5px 0;
  display: flex;
  justify-content: center;
}

/* prettier-ignore */
@media (max-width: 1200px) { #photos { column-count: 3; } }
/* prettier-ignore */
@media (max-width: 1000px) { #photos { column-count: 2; } }
/* prettier-ignore */
@media (max-width: 800px)  { #photos { column-count: 1; } }

.sidebar-slide-enter-active,
.sidebar-slide-leave-active {
  transition: transform 0.2s ease;
}

.sidebar-slide-enter-from,
.sidebar-slide-leave-to {
  transform: translateX(100%);
  transition: all 150ms ease-in 0s;
}
</style>
