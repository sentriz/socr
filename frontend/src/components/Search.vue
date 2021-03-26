<template>
  <div class="space-y-6">
    <div class="flex flex-col md:flex-row gap-2">
      <input v-model="reqQuery" class="inp w-full" type="text" placeholder="enter screenshot text query" />
      <SearchFilter label="sort by" :items="reqSortOptions" v-model:selected="reqSortOption" />
      <SearchFilter label="directory" :items="reqFilterOptions" v-model:selected="reqFilterOption" />
    </div>
    <div ref="scroller">
      <p v-if="!loading" class="text-gray-500 text-right">fetched {{ respTook.toFixed(2) }}ms</p>
      <div v-for="(page, i) in respPages" :key="i" class="mt-2">
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
  <UploadHint />
</template>

<script setup lang="ts">
import ScreenshotHighlight from './ScreenshotHighlight.vue'
import ScreenshotBackground from './ScreenshotBackground.vue'
import SearchSidebar from './SearchSidebar.vue'
import SearchFilter from './SearchFilter.vue'
import LoadingSpinner from './LoadingSpinner.vue'
import ClipboardUploader from './ClipboardUploader.vue'
import UploadHint from './UploadHint.vue'

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

const reqQuery = ref('')
const reqQueryDebounced = useDebounce(reqQuery, 100)
const reqPageSize = 25
const reqPageNum = ref(0)

type Sort = { label: string; icon: string; field: string; order: SortOrder }
const reqSortSimilarity: Sort = { label: 'similarity', icon: 'fas fa-search', field: 'similarity', order: SortOrder.Desc }
const reqSortTimeDesc: Sort = { label: 'time desc', icon: 'fas fa-chevron-down', field: 'timestamp', order: SortOrder.Desc }
const reqSortTimeAsc: Sort = { label: 'time asc', icon: 'fas fa-chevron-up', field: 'timestamp', order: SortOrder.Asc }
const reqSortOptions = ref([reqSortTimeDesc, reqSortTimeAsc])
const reqSortOption = ref(reqSortTimeDesc)

type Filter = { label: string; icon: string; directory?: string }
const reqFilterAll: Filter = { label: 'all', icon: 'fas fa-asterisk' }
const reqFilterOptions = ref([reqFilterAll])
const reqFilterOption = ref(reqFilterAll)

const respTook = ref(0)
const respHasMore = ref(true)
const respPages = ref<string[][]>([])

const fetchScreenshots = async () => {
  if (loading.value) return
  if (!respHasMore.value) return

  console.log('loading page #%d', reqPageNum.value)
  const offset = reqPageSize * reqPageNum.value
  const sort = { field: reqSortOption.value.field, order: reqSortOption.value.order }
  const resp = await load(reqPageSize, offset, sort, reqQuery.value, reqFilterOption.value.directory)
  if (isError(resp)) return

  respTook.value = (resp.result.took || 0) / 10 ** 6
  respHasMore.value = !!resp.result.screenshots?.length
  if (!respHasMore.value) return

  reqPageNum.value++
  respPages.value.push([])
  for (const screenshot of resp.result.screenshots || []) {
    respPages.value[respPages.value.length - 1].push(screenshot.hash)
  }
}

const clearParameters = () => {
  reqPageNum.value = 0
  respPages.value = []
  respHasMore.value = true

  // 🤔
  if (reqQuery.value && !reqSortOptions.value.includes(reqSortSimilarity)) {
    reqSortOptions.value.splice(0, 0, reqSortSimilarity)
    reqSortOption.value = reqSortSimilarity
  }
  if (!reqQuery.value && reqSortOptions.value.includes(reqSortSimilarity)) {
    reqSortOptions.value.splice(0, 1)
    reqSortOption.value = reqSortTimeDesc
  }
}

const clearFetchScreenshots = () => {
  clearParameters()
  fetchScreenshots()
}

// fetch screenshots on reaching the bottom of page
const scroller = useInfiniteScroll(fetchScreenshots)

// fetch screenshots on filter, sort, or query change
watch(reqSortOption, clearFetchScreenshots)
watch(reqFilterOption, clearFetchScreenshots)
watch(reqQueryDebounced, clearFetchScreenshots, { immediate: true })

// fetch directories for dropdown on mount
onMounted(async () => {
  const resp = await reqDirectories()
  if (isError(resp)) return

  reqFilterOptions.value = [
    ...reqFilterOptions.value,
    ...resp.result.map((d) => ({
      label: d.directory_alias,
      icon: d.is_uploads ? 'fas fa-folder-plus' : 'fas fa-folder',
      directory: d.directory_alias,
    })),
  ]
})
</script>
