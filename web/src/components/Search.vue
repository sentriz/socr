<template>
  <div class="space-y-6">
    <div class="md:grid-cols-2 lg:grid-cols-12 grid grid-cols-1 gap-2">
      <input class="lg:col-span-9 inp" v-model="reqQuery" type="text" placeholder="enter media text query" />
      <search-filter class="lg:col-span-3" label="sort by" :items="reqSortOptions" v-model:selected="reqSort" />
      <search-filter class="lg:col-span-3" label="media" :items="reqMediaOptions" v-model:selected="reqMedia" />
      <search-filter class="lg:col-span-3" label="directory" :items="reqDirOptions" v-model:selected="reqDir" />
      <search-filter class="lg:col-span-3" label="year" :items="reqYearOptions" v-model:selected="reqYear" />
      <search-filter class="lg:col-span-3" label="month" :items="reqMonthOptions" v-model:selected="reqMonth" :disabled="!reqYear.year" />
    </div>
    <p v-if="!loading" class="text-right text-gray-500">fetched {{ respTook.toFixed(2) }}ms</p>
    <div v-for="(page, i) in respPages" :key="i">
      <div v-show="i !== 0" class="my-6">
        <span class="text-gray-500"> page {{ i + 1 }}</span>
        <hr class="m-0" />
      </div>
      <div class="col-resp col-gap-4 space-y-4">
        <media-background v-for="hash in page" :key="hash" :hash="hash" class="shadow-lg max-h-[600px] overflow-y-auto flex justify-center">
          <router-link :to="{ name: 'search', params: { hash } }" class="block">
            <media-highlight thumb :hash="hash" />
          </router-link>
        </media-background>
      </div>
    </div>
    <loading-spinner v-if="loading" />
    <search-no-results v-else-if="respPages.length === 0" />
  </div>
  <div ref="scroller" />
  <teleport to="#overlays"><search-sidebar :hash="sidebarHash" /></teleport>
  <teleport to="#overlays"><uploader-clipboard /></teleport>
  <teleport to="#overlays"><uploader-file /></teleport>
</template>

<script setup lang="ts">
import MediaHighlight from './MediaHighlight.vue'
import MediaBackground from './MediaBackground.vue'
import SearchSidebar from './SearchSidebar.vue'
import SearchFilter from './SearchFilter.vue'
import LoadingSpinner from './LoadingSpinner.vue'
import UploaderClipboard from './UploaderClipboard.vue'
import UploaderFile from './UploaderFile.vue'
import SearchNoResults from './SearchNoResults.vue'

import { watch, computed, onMounted, ref } from 'vue'
import type { Component } from 'vue'
import { useRoute } from 'vue-router'
import { useDebounce } from '@vueuse/core'
import { isError, SortOrder, reqDirectories, MediaType } from '../api'
import type { PayloadSearch } from '../api'
import useStore from '../composables/useStore'
import useInfiniteScroll from '../composables/useInfiniteScroll'
import useLoading from '../composables/useLoading'
import {
  SearchIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  DocumentDuplicateIcon,
  CameraIcon,
  VideoCameraIcon,
  FolderAddIcon,
  FolderIcon,
  CalendarIcon,
  GlobeIcon,
} from '@heroicons/vue/outline'
import router from '../router'

const store = useStore()
const route = useRoute()
const { loading, load } = useLoading(store.loadMedias)

const sidebarHash = computed(() => (route.params.hash as string) || '')

const reqQuery = ref('')
const reqPageNum = ref(0)
const reqQueryDebounced = useDebounce(reqQuery, 500)
const reqPageSize = 25

type Sort = { label: string; icon: Component; field: string; order: SortOrder }
const reqSortSimilarity: Sort = { label: 'similarity', icon: SearchIcon, field: 'similarity', order: SortOrder.Desc }
const reqSortTimeDesc: Sort = { label: 'time desc', icon: ChevronDownIcon, field: 'timestamp', order: SortOrder.Desc }
const reqSortTimeAsc: Sort = { label: 'time asc', icon: ChevronUpIcon, field: 'timestamp', order: SortOrder.Asc }
const reqSortOptionsDefault = [reqSortTimeDesc, reqSortTimeAsc]
const reqSortOptionsSimilarity = [reqSortSimilarity, reqSortTimeDesc, reqSortTimeAsc]
const reqSortOptions = ref(reqSortOptionsDefault)
const reqSort = ref(reqSortTimeDesc)

type Dir = { label: string; icon: Component; directory?: string }
const reqDirAll: Dir = { label: 'all', icon: DocumentDuplicateIcon }
const reqDirOptions = ref([reqDirAll])
const reqDir = ref(reqDirAll)

onMounted(async () => {
  const resp = await reqDirectories()
  if (isError(resp)) return

  for (const d of resp.result)
    reqDirOptions.value.push({
      label: d.directory_alias,
      icon: d.is_uploads ? FolderAddIcon : FolderIcon,
      directory: d.directory_alias,
    })
})

type Media = { label: string; icon: Component; media?: MediaType }
const reqMediaAny: Media = { label: 'any', icon: DocumentDuplicateIcon }
const reqMediaImage: Media = { label: 'image', icon: CameraIcon, media: MediaType.Image }
const reqMediaVideo: Media = { label: 'video', icon: VideoCameraIcon, media: MediaType.Video }
const reqMediaOptions = ref([reqMediaAny, reqMediaImage, reqMediaVideo])
const reqMedia = ref(reqMediaAny)

const currentYear = new Date().getFullYear()
type Year = { label: string; icon: Component; year?: number }
const reqYearAny: Year = { label: `any`, icon: CalendarIcon }
const reqYearOptions = ref([reqYearAny])
const reqYear = ref(reqYearAny)

for (let i = currentYear; i >= currentYear - 10; i--) {
  reqYearOptions.value.push({ label: `${i}`, icon: CalendarIcon, year: i })
}

type Month = { label: string; icon: Component; month?: number }
const reqMonthAny: Month = { label: `any`, icon: GlobeIcon }
const reqMonthOptions = ref([reqMonthAny])
const reqMonth = ref(reqMonthAny)

const months = ['jan', 'feb', 'mar', 'apr', 'may', 'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec']
for (const [month, label] of months.entries()) {
  reqMonthOptions.value.push({ label, icon: CalendarIcon, month })
}

const respTook = ref(0)
const respHasMore = ref(true)
const respPages = ref<string[][]>([])

const fetchMedias = async () => {
  if (loading.value) return
  if (!respHasMore.value) return

  const req: PayloadSearch = {
    body: reqQuery.value,
    limit: reqPageSize,
    offset: reqPageSize * reqPageNum.value,
    sort: { field: reqSort.value.field, order: reqSort.value.order },
    directory: reqDir.value.directory,
    media: reqMedia.value.media,
  }
  if (reqYear.value.year) {
    req.date_from = new Date(reqYear.value.year, reqMonth.value.month ?? 0, 1)
    req.date_to = new Date(reqYear.value.year, reqMonth.value.month ?? 11, 31)
  }

  console.log('loading page #%d', reqPageNum.value)
  const resp = await load(req)
  if (isError(resp)) return

  respTook.value = (resp.result.took || 0) / 10 ** 6
  respHasMore.value = !!resp.result.medias?.length
  if (!respHasMore.value) return

  reqPageNum.value++
  respPages.value.push([])
  for (const media of resp.result.medias || []) {
    respPages.value[respPages.value.length - 1].push(media.hash)
  }
}

const resetParameters = () => {
  reqPageNum.value = 0
  respPages.value = []
  respHasMore.value = true
}

const resetSortOptions = () => {
  if (reqQuery.value && !reqSortOptions.value.includes(reqSortSimilarity)) {
    reqSortOptions.value = reqSortOptionsSimilarity
    reqSort.value = reqSortSimilarity
  }
  if (!reqQuery.value && reqSortOptions.value.includes(reqSortSimilarity)) {
    reqSortOptions.value = reqSortOptionsDefault
    reqSort.value = reqSortTimeDesc
  }
}

const scroller = useInfiniteScroll(fetchMedias)

watch([reqSort, reqDir, reqMedia, reqYear, reqMonth, reqQueryDebounced], () => {
  resetParameters()
  resetSortOptions()
  fetchMedias()
})

onMounted(async () => {
  await fetchMedias()
})

document.onkeydown = (e: KeyboardEvent) => {
  let dir = 0
  if (e.key === 'ArrowRight') dir = 1
  else if (e.key === 'ArrowLeft') dir = -1
  else return

  if (!sidebarHash.value) return
  if (!respPages.value?.length) return

  const flat = respPages.value.flat()
  const idx = flat.findIndex((e) => e === sidebarHash.value)
  const hash = flat[(idx + dir) % flat.length]
  router.replace({ name: 'search', params: { hash } })
}
</script>
