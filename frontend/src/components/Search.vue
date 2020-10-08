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
  <div ref="intersectionTrigger"></div>
  <SearchSidebar :id="sidebarID" />
</template>

<script setup="props">
import ScreenshotHighlight from "./ScreenshotHighlight.vue";
import SearchSidebar from "./SearchSidebar.vue";
export default {
  components: { ScreenshotHighlight, SearchSidebar },
  props: {},
};

import { inject, ref, reactive, watch, computed } from "vue";
import { useRoute } from "vue-router";
import { useThrottle } from "@vueuse/core";
import { makeUseInfiniteScroll } from "vue-use-infinite-scroll";
import { reqSearch, reqSearchParams } from "../api";

const useInfiniteScroll = makeUseInfiniteScroll({});
export const intersectionTrigger = ref(null);

const pageSize = 25;
const pageNum = useInfiniteScroll(intersectionTrigger);

watch(pageNum, () => {
  console.log("new page %d", pageNum.value);
  fetchScreenshots();
});

export const query = ref("");
export const queryThrottled = useThrottle(query, 250);
watch(queryThrottled, () => {
  pageNum.value = 0;
  store.screenshots = {};
  fetchScreenshots();
});

const route = useRoute();
export const sidebarID = computed(() => route.params.id);

export const store = inject("store");
export const reqTotalHits = ref(0);
export const reqTookMs = ref(0);
export const fetchScreenshots = async () => {
  const resp = await reqSearch(
    reqSearchParams(pageSize, pageSize * pageNum.value, query.value),
  );

  reqTotalHits.value = resp.total_hits;
  reqTookMs.value = Math.round((resp.took / 100000) * 100) / 100;
  for (const hit of resp.hits) {
    store.screenshots[hit.id] = hit;
  }
};
</script>
