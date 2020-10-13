<template>
  <input
    v-model="reqQuery"
    class="inp w-full"
    type="text"
    placeholder="enter screenshot text query"
  />
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
  <div v-if="isLoading" class="bg-gray-300 text-gray-600 text-center rounded p-3 my-6">
    <i class="animate-spin fas fa-circle-notch"></i> loading more
  </div>
  <SearchSidebar :id="sidebarID" />
</template>

<script setup="props">
import ScreenshotHighlight from "./ScreenshotHighlight.vue";
import SearchSidebar from "./SearchSidebar.vue";
export default {
  components: { ScreenshotHighlight, SearchSidebar },
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
export const pages = ref([]);

export const reqQuery = ref("");

const reqParamSortingMode = "timestamp_desc";
const reqParamSorting = {
  timestamp_desc: [`-${fields.TIMESTAMP}`],
  timestamp_asc: [`${fields.TIMESTAMP}`],
};

export const reqTotalHits = ref(0);
export const reqTookMs = ref(0);
export const fetchScreenshots = async (pageNum) => {
  console.log("loading page #%d", pageNum);

  const from = pageSize * pageNum;
  const sort = reqParamSorting[reqParamSortingMode];
  const resp = reqQuery.value
    ? await store.screenshotsLoadTerm(pageSize, from, sort, reqQuery.value)
    : await store.screenshotsLoadRecent(pageSize, from, sort);

  pages.value.push([]);
  for (const hit of resp.hits) {
    pages.value[pages.value.length - 1].push(hit.id);
  }

  reqTotalHits.value = resp.total_hits;
  reqTookMs.value = Math.round((resp.took / 100000) * 100) / 100;
};

export const { scroller, pageNum, isLoading } = useInfiniteScroll(fetchScreenshots);

export const reqQueryThrottled = useThrottle(reqQuery, 250);
watch(
  reqQueryThrottled,
  () => {
    pages.value = [];
    fetchScreenshots(0);
    pageNum.value = 1;
  },
  { immediate: true },
);
</script>
