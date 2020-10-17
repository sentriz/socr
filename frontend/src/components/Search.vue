<template>
  <div class="block space-y-2 sm:flex sm:space-y-0 sm:space-x-2">
    <input
      v-model="reqQuery"
      class="inp w-full"
      type="text"
      placeholder="enter screenshot text query"
    />
    <SearchSortFilter :items="reqParamSortModes" v-model="reqParamSortMode" />
  </div>
  <p class="my-3 text-gray-500 text-right">
    {{ reqTotalHits }} results found in {{ reqTookMs }}ms
  </p>
  <div ref="scroller">
    <div v-for="(page, i) in pages">
      <div v-show="i !== 0" class="my-6">
        <span class="text-gray-500"> page {{ i + 1 }}</span>
        <hr class="m-0" />
      </div>
      <div class="col-resp gap-x-3 space-y-3">
        <div v-for="screenshotID in page">
          <router-link :to="{ name: 'search', params: { id: screenshotID } }">
            <ScreenshotHighlight
              :id="screenshotID"
              class="border border-gray-300 rounded"
            />
          </router-link>
        </div>
      </div>
    </div>
  </div>
  <div v-if="isLoading" class="bg-gray-300 text-gray-600 text-center rounded p-3 mt-6">
    <i class="animate-spin fas fa-circle-notch"></i> loading more
  </div>
  <SearchSidebar :id="sidebarID" />
</template>

<script setup="props">
import ScreenshotHighlight from "./ScreenshotHighlight.vue";
import SearchSidebar from "./SearchSidebar.vue";
import SearchSortFilter from "./SearchSortFilter.vue";
export default {
  components: { ScreenshotHighlight, SearchSidebar, SearchSortFilter },
  props: {},
};

import { inject, ref, reactive, watch, computed, onMounted, onUnmounted } from "vue";
import { useRoute } from "vue-router";
import { useThrottle } from "@vueuse/core";
import { fields } from "../api";
import { useStore } from "../store";
import useInfiniteScroll from "../composables/useInfiniteScroll";

export const store = useStore();
export const route = useRoute();

export const sidebarID = computed(() => route.params.id);

const pageSize = 25;
const pageNum = ref(0);
export const pages = ref([]);

export const reqParamSortMode = ref(0);
export const reqParamSortModes = [
  { filter: [`-${fields.TIMESTAMP}`], name: "updated", icon: "fa-chevron-down" },
  { filter: [`${fields.TIMESTAMP}`], name: "updated", icon: "fa-chevron-up" },
];

export const reqQuery = ref("");
export const reqQueryThrottled = useThrottle(reqQuery, 250);
export const reqTotalHits = ref(0);
export const reqTookMs = ref(0);
const fetchScreenshots = async () => {
  console.log("loading page #%d", pageNum.value);

  const from = pageSize * pageNum.value;
  const sort = reqParamSortModes[reqParamSortMode.value].filter;
  const resp = await store.screenshotsLoad(pageSize, from, sort, reqQuery.value);

  pages.value.push([]);
  for (const hit of resp.hits) {
    pages.value[pages.value.length - 1].push(hit.id);
  }

  reqTotalHits.value = resp.total_hits;
  reqTookMs.value = Math.round((resp.took / 1000000) * 100) / 100;
};

const fetchScreenshotsClear = async () => {
  pageNum.value = 0;
  pages.value = [];
  await fetchScreenshots();
  pageNum.value++;
};

// fetch screenshots on filter, sort, and mount
watch(reqParamSortMode, fetchScreenshotsClear);
watch(reqQueryThrottled, fetchScreenshotsClear, { immediate: true });

// fetch screenshots on reaching the bottom of page
export const { scroller, isLoading } = useInfiniteScroll(async () => {
  await fetchScreenshots();
  pageNum.value++;
});
</script>
