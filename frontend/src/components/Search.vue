<template>
  <div class="space-y-6">
    <div class="flex flex-col md:flex-row gap-2">
      <input v-model="reqQuery" class="inp w-full" type="text" placeholder="enter screenshot text query" />
      <SearchFilter label="sort by" :items="reqSortOptions" v-model:selected="reqSortOption" />
      <SearchFilter label="directory" :items="reqFilterOptions" v-model:selected="reqFilterOption" />
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
    <LoadingSpinner v-if="loading" />
  </div>
  <SearchSidebar :hash="sidebarHash" />
  <ClipboardUploader />
</template>

<script setup lang="ts">
import ScreenshotHighlight from './ScreenshotHighlight.vue'
import ScreenshotBackground from './ScreenshotBackground.vue'
import SearchSidebar from './SearchSidebar.vue'
import SearchFilter from './SearchFilter.vue'
import LoadingSpinner from './LoadingSpinner.vue'
import ClipboardUploader from './ClipboardUploader.vue'

import { ref, watch, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useDebounce } from '@vueuse/core'
import { isError, SortOrder, reqDirectories } from '../api'
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

type FilterDisplayRow = { label: string; icon: string }
type SortRow = FilterDisplayRow & { field: string; order: SortOrder }
type FilterRow = FilterDisplayRow & { directory?: string }

const reqSortOptions: SortRow[] = [
  { label: 'time desc', icon: 'fas fa-chevron-down', field: 'timestamp', order: SortOrder.Desc },
  { label: 'time asc', icon: 'fas fa-chevron-up', field: 'timestamp', order: SortOrder.Asc },
]
const reqSortOption = ref<SortRow>(reqSortOptions[0])

const reqFilterOptions = ref<FilterRow[]>([{ label: 'all', icon: 'fas fa-asterisk' }])
const reqFilterOption = ref<FilterRow>(reqFilterOptions.value[0])

onMounted(async () => {
  const resp = await reqDirectories()
  if (isError(resp)) return

  for (const d of resp.result)
    reqFilterOptions.value.push({
      label: d.directory_alias,
      icon: d.is_uploads ? 'fas fa-folder-plus' : 'fas fa-folder',
      directory: d.directory_alias,
    })
})

const reqQuery = ref('')
const reqQueryDebounced = useDebounce(reqQuery, 100)

const respTook = ref(0)
const respHasMore = ref(true)

const fetchScreenshots = async () => {
  if (loading.value) return
  if (!respHasMore.value) return

  console.log('loading page #%d', pageNum.value)
  const offset = pageSize * pageNum.value
  const sort = { field: reqSortOption.value.field, order: reqSortOption.value.order }
  const resp = await load(pageSize, offset, sort, reqQuery.value, reqFilterOption.value.directory)
  if (isError(resp)) return

  respTook.value = (resp.result.took || 0) / 10 ** 6
  respHasMore.value = !!resp.result.screenshots?.length
  if (!respHasMore.value) return

  pageNum.value++
  pages.value.push([])
  for (const screenshot of resp.result.screenshots || []) {
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
watch(reqSortOption, fetchScreenshotsClear)
watch(reqFilterOption, fetchScreenshotsClear)
watch(reqQueryDebounced, fetchScreenshotsClear, { immediate: true })

// fetch screenshots on reaching the bottom of page
const scroller = useInfiniteScroll(fetchScreenshots)
</script>
