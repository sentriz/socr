<template>
  <div class="space-y-6">
    <div class="block space-y-2 sm:flex sm:space-y-0 sm:space-x-2">
      <input
        v-model="reqQuery"
        class="inp w-full"
        type="text"
        placeholder="enter screenshot text query"
      />
      <SearchSortFilter :items="reqParamSortModes" v-model="reqParamSortMode" />
    </div>
    <div ref="scroller">
      <p v-if="!reqIsLoading" class="text-gray-500 text-right">
        {{ reqTotalHits }} results found in {{ reqTookMs }}ms
      </p>
      <div v-for="(page, i) in pages" class="mt-1">
        <div v-show="i !== 0" class="my-6">
          <span class="text-gray-500"> page {{ i + 1 }}</span>
          <hr class="m-0" />
        </div>
        <div class="col-resp gap-x-4 space-y-4">
          <div v-for="screenshotID in page" class="bg-gray-200">
            <router-link :to="{ name: 'search', params: { id: screenshotID } }">
              <ScreenshotHighlight :id="screenshotID" class="mx-auto" />
            </router-link>
          </div>
        </div>
      </div>
    </div>
    <div v-if="reqIsLoading" class="bg-gray-300 text-gray-600 text-center rounded p-3">
      <i class="animate-spin fas fa-circle-notch"></i> loading
    </div>
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
import { useDebounce } from "@vueuse/core";
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
  { filter: [`-${fields.TIMESTAMP}`], name: "date", icon: "fa-chevron-down" },
  { filter: [`${fields.TIMESTAMP}`], name: "date", icon: "fa-chevron-up" },
];

export const reqQuery = ref("");
export const reqQueryDebounced = useDebounce(reqQuery, 500);
export const reqTotalHits = ref(0);
export const reqTookMs = ref(0);
export const reqHasMore = ref(true);
export const reqIsLoading = ref(false);
const fetchScreenshots = async () => {
  if (reqIsLoading.value) return;
  if (!reqHasMore.value) return;

  console.log("loading page #%d", pageNum.value);

  const from = pageSize * pageNum.value;
  const sort = reqParamSortModes[reqParamSortMode.value].filter;

  reqIsLoading.value = true;
  const resp = await store.screenshotsLoad(pageSize, from, sort, reqQuery.value);
  reqIsLoading.value = false;

  reqTotalHits.value = resp.total_hits;
  reqTookMs.value = Math.round((resp.took / 1000000) * 100) / 100;
  reqHasMore.value = from + resp.hits.length < resp.total_hits;

  pages.value.push([]);
  for (const hit of resp.hits) {
    pages.value[pages.value.length - 1].push(hit.id);
  }

  pageNum.value++;
};

const fetchScreenshotsClear = async () => {
  reqHasMore.value = true;
  pageNum.value = 0;
  pages.value = [];
  await fetchScreenshots();
};

// fetch screenshots on filter, sort, and mount
watch(reqParamSortMode, fetchScreenshotsClear);
watch(reqQueryDebounced, fetchScreenshotsClear, { immediate: true });

// fetch screenshots on reaching the bottom of page
export const scroller = useInfiniteScroll(fetchScreenshots);
</script>
