<template>
  <div class="space-y-6">
    <div class="flex flex-col md:flex-row gap-2">
      <input v-model="reqQuery" class="inp w-full" type="text" placeholder="enter screenshot text query" />
      <SearchSortFilter v-if="!reqQuery" v-model:field="reqSortField" v-model:order="reqSortOrder" label="date" />
    </div>
    <div ref="scroller">
      <p v-if="!loading" class="text-gray-500 text-right">fetched {{ respTook.toFixed(2) }}ms</p>
      <div v-for="(page, i) in pages" :key="i" class="mt-2">
        <div v-show="i !== 0" class="my-6">
          <span class="text-gray-500"> page {{ i + 1 }}</span>
          <hr class="m-0" />
        </div>
        <div class="col-resp gap-x-4 space-y-4">
          <ScreenshotBackground v-for="hash in page" :key="hash" :hash="hash" class="shadow-lg">
            <router-link :to="{ name: 'search', params: { hash } }">
              <ScreenshotHighlight :hash="hash" class="mx-auto" />
            </router-link>
          </ScreenshotBackground>
        </div>
      </div>
    </div>
    <SearchLoading v-if="loading" />
  </div>
  <SearchSidebar :hash="sidebarHash" />
  <ClipboardUploader />
</template>

<script setup lang="ts">
import ScreenshotHighlight from './ScreenshotHighlight.vue'
import ScreenshotBackground from './ScreenshotBackground.vue'
import SearchSidebar from './SearchSidebar.vue'
import SearchSortFilter from './SearchSortFilter.vue'
import SearchLoading from './SearchLoading.vue'
import ClipboardUploader from './ClipboardUploader.vue'

import { ref, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useDebounce } from '@vueuse/core'
import { isError, SortOrder } from '../api'
import useStore from '../composables/useStore'
import useInfiniteScroll from '../composables/useInfiniteScroll'
import useLoading from '../composables/useLoading'

const store = useStore()
const route = useRoute()
const { loading, load } = useLoading(store.loadScreenshots)

const sidebarHash = computed(() => (route.params.hash as string) || '')

const pageSize = 25
const pageNum = ref(0)
const pages = ref<string[][]>([])

const reqSortField = ref('timestamp')
const reqSortOrder = ref(SortOrder.Desc)
const reqQuery = ref('')
const reqQueryDebounced = useDebounce(reqQuery, 100)

const respTook = ref(0)
const respHasMore = ref(true)

const fetchScreenshots = async () => {
  if (loading.value) return
  if (!respHasMore.value) return

  console.log('loading page #%d', pageNum.value)
  const from = pageSize * pageNum.value
  const sort = { field: reqSortField.value, order: reqSortOrder.value }
  const resp = await load(pageSize, from, sort, reqQuery.value)
  if (isError(resp)) return

  respTook.value = (resp.result.took || 0) / 10 ** 6
  respHasMore.value = !!resp.result.screenshots.length
  if (!respHasMore.value) return

  pageNum.value++
  pages.value.push([])
  for (const screenshot of resp.result.screenshots) {
    pages.value[pages.value.length - 1].push(screenshot.hash)
  }
}

const fetchScreenshotsClear = async () => {
  pageNum.value = 0
  pages.value = []
  respHasMore.value = true
  await fetchScreenshots()
}

// fetch screenshots on filter, sort, and mount
watch(reqSortOrder, fetchScreenshotsClear)
watch(reqQueryDebounced, fetchScreenshotsClear, { immediate: true })

// fetch screenshots on reaching the bottom of page
const scroller = useInfiniteScroll(fetchScreenshots)
</script>
