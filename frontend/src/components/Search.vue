<template>
  <div class="space-y-6">
    <div class="flex flex-col md:flex-row gap-2">
      <input v-model="reqQuery" class="inp w-full" type="text" placeholder="enter screenshot text query" />
      <SearchSortFilter
        v-model:field="reqSortField"
        v-model:order="reqSortOrder"
        label="date"
      />
    </div>
    <div ref="scroller">
      <p v-if="!loading" class="text-gray-500 text-right">
        {{ respLength }} results found in {{ respTookMs || 0 }}ms
      </p>
      <div v-for="(page, i) in pages" :key="i" class="mt-2">
        <div v-show="i !== 0" class="my-6">
          <span class="text-gray-500"> page {{ i + 1 }}</span>
          <hr class="m-0" />
        </div>
        <div class="col-resp gap-x-4 space-y-4">
          <ScreenshotBackground v-for="hash in page" :key="hash" :hash="hash" class="shadow-lg">
            <router-link :to="{ name: 'search', params: { hash } }">
              <ScreenshotHighlight :hash="hash" class="mx-auto"/>
            </router-link>
          </ScreenshotBackground>
        </div>
      </div>
    </div>
    <SearchLoading v-if="loading" />
  </div>
  <SearchSidebar :id="sidebarID" />
</template>

<script setup lang="ts">
import ScreenshotHighlight from "./ScreenshotHighlight.vue";
import ScreenshotBackground from "./ScreenshotBackground.vue";
import SearchSidebar from "./SearchSidebar.vue";
import SearchSortFilter from "./SearchSortFilter.vue";
import SearchLoading from './SearchLoading.vue'

import { ref, watch, computed } from "vue";
import { useRoute } from "vue-router";
import { useDebounce } from "@vueuse/core";
import { isError, SortOrder } from "../api";
import type { Search } from "../api";
import useStore from "../composables/useStore";
import useInfiniteScroll from "../composables/useInfiniteScroll";
import useLoading from "../composables/useLoading";

const store = useStore();
const route = useRoute();
const { loading, load } = useLoading(store.loadScreenshots);

const sidebarID = computed(() => route.params.id as string || "");

const pageSize = 25;
const pageNum = ref(0);
const pages = ref<string[][]>([]);

const reqSortField = ref("timestamp")
const reqSortOrder = ref(SortOrder.Desc)
const reqQuery = ref("");
const reqQueryDebounced = useDebounce(reqQuery, 500);

const resp = ref<Search>();
const hasMore = ref(true);

const fetchScreenshots = async () => {
  if (loading.value) return;
  if (!hasMore.value) return;

  console.log("loading page #%d", pageNum.value);
  const from = pageSize * pageNum.value;
  const sort = { field: reqSortField.value, order: reqSortOrder.value }
  const r = await load(pageSize, from, sort, reqQuery.value);
  if (isError(r)) return

  resp.value = r.result
  hasMore.value = from + resp.value.length < resp.value.total;
  pageNum.value++;
  pages.value.push([]);
  for (const screenshot of resp.value.screenshots) {
    pages.value[pages.value.length - 1].push(screenshot.hash);
  }
};

const fetchScreenshotsClear = async () => {
  hasMore.value = true;
  pageNum.value = 0;
  pages.value = [];
  await fetchScreenshots();
};

const respLength = computed(() => resp.value?.length || 0);
const respTookMs = computed(() => ((resp.value?.took || 0) / 10 ** 6).toFixed(2))

// fetch screenshots on filter, sort, and mount
watch(reqSortOrder, fetchScreenshotsClear);
watch(reqQueryDebounced, fetchScreenshotsClear, { immediate: true });

// fetch screenshots on reaching the bottom of page
const scroller = useInfiniteScroll(fetchScreenshots);
</script>
