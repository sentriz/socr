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
      <router-link :to="{ name: 'search', params: { id: screenshot.id } }">
        <ScreenshotHighlight :id="screenshot.id" class="border border-gray-300 rounded" />
      </router-link>
    </div>
  </div>
  <SearchSidebar :id="sidebarID" />
</template>

<script setup>
import ScreenshotHighlight from "./ScreenshotHighlight.vue";
import SearchSidebar from "./SearchSidebar.vue";
export default {
  props: {},
  components: {
    ScreenshotHighlight,
    SearchSidebar,
  },
};

import { inject, ref, reactive, watch, computed } from "vue";
import { useRoute } from "vue-router";
import { useThrottle } from "@vueuse/core";
import { reqSearch, reqSearchParams } from "../api";

export const query = ref("");
export const queryThrottled = useThrottle(query, 250);
watch(queryThrottled, (v, _) => {
  store.screenshots = {};
  if (v) fetchScreenshots();
});

const route = useRoute();
export const sidebarID = computed(() => route.params.id);

export const store = inject("store");
export const reqTotalHits = ref(0);
export const reqTookMs = ref(0);
export const fetchScreenshots = async () => {
  const resp = await reqSearch(reqSearchParams(query.value));

  reqTotalHits.value = resp.total_hits;
  reqTookMs.value = Math.round((resp.took / 100000) * 100) / 100;
  for (const hit of resp.hits) {
    store.screenshots[hit.id] = hit;
  }
};
</script>
