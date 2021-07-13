<template>
  <div class="space-y-6">
    <div class="md:flex-row flex flex-col gap-2">
      <input v-model="reqQuery" class="inp w-full" type="text" placeholder="enter media text query" />
      <search-filter label="sort by" :items="reqSortOptions" v-model:selected="reqSortOption" />
      <search-filter label="media" :items="reqMediaOptions" v-model:selected="reqMediaOption" />
      <search-filter label="directory" :items="reqDirOptions" v-model:selected="reqDirOption" />
    </div>
    <div ref="scroller">
      <p v-if="!loading" class="text-right text-gray-500">fetched {{ respTook.toFixed(2) }}ms</p>
      <div v-for="(page, i) in respPages" :key="i" class="mt-2">
        <div v-show="i !== 0" class="my-6">
          <span class="text-gray-500"> page {{ i + 1 }}</span>
          <hr class="m-0" />
        </div>
        <div class="col-resp gap-x-4 space-y-4">
          <media-background v-for="hash in page" :key="hash" :hash="hash" class="shadow-lg">
            <router-link :to="{ name: 'search', params: { hash } }">
              <media-highlight :hash="hash" class="mx-auto" />
            </router-link>
          </media-background>
        </div>
      </div>
    </div>
    <loading-spinner v-if="loading" />
    <search-no-results v-else-if="respPages.length === 0" />
  </div>
  <search-sidebar :hash="sidebarHash" />
  <uploader-clipboard />
  <uploader-file />
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
} from 'heroicons-vue3/outline'

const store = useStore()
const route = useRoute()
const { loading, load } = useLoading(store.loadMedias)

const sidebarHash = computed(() => (route.params.hash as string) || '')

const reqQuery = ref('')
const reqPageNum = ref(0)
const reqQueryDebounced = useDebounce(reqQuery, 200)
const reqPageSize = 25

type Sort = { label: string; icon: Component; field: string; order: SortOrder }
const reqSortSimilarity: Sort = { label: 'similarity', icon: SearchIcon, field: 'similarity', order: SortOrder.Desc }
const reqSortTimeDesc: Sort = { label: 'time desc', icon: ChevronDownIcon, field: 'timestamp', order: SortOrder.Desc }
const reqSortTimeAsc: Sort = { label: 'time asc', icon: ChevronUpIcon, field: 'timestamp', order: SortOrder.Asc }
const reqSortOptionsDefault = [reqSortTimeDesc, reqSortTimeAsc]
const reqSortOptionsSimilarity = [reqSortSimilarity, reqSortTimeDesc, reqSortTimeAsc]
const reqSortOptions = ref(reqSortOptionsDefault)
const reqSortOption = ref(reqSortTimeDesc)

type Dir = { label: string; icon: Component; directory?: string }
const reqDirAll: Dir = { label: 'all', icon: DocumentDuplicateIcon }
const reqDirOptions = ref([reqDirAll])
const reqDirOption = ref(reqDirAll)

type Media = { label: string; icon: Component; media?: MediaType }
const reqMediaAny: Media = { label: 'any', icon: DocumentDuplicateIcon }
const reqMediaImage: Media = { label: 'image', icon: CameraIcon, media: MediaType.Image }
const reqMediaVideo: Media = { label: 'video', icon: VideoCameraIcon, media: MediaType.Video }
const reqMediaOptions = ref([reqMediaAny, reqMediaImage, reqMediaVideo])
const reqMediaOption = ref(reqMediaAny)

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
    sort: { field: reqSortOption.value.field, order: reqSortOption.value.order },
    directory: reqDirOption.value.directory,
    media: reqMediaOption.value.media,
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

const resetDirs = () => {
  if (reqQuery.value && !reqSortOptions.value.includes(reqSortSimilarity)) {
    reqSortOptions.value = reqSortOptionsSimilarity
    reqSortOption.value = reqSortSimilarity
  }
  if (!reqQuery.value && reqSortOptions.value.includes(reqSortSimilarity)) {
    reqSortOptions.value = reqSortOptionsDefault
    reqSortOption.value = reqSortTimeDesc
  }
}

const scroller = useInfiniteScroll(fetchMedias)

watch([reqSortOption, reqDirOption, reqMediaOption, reqQueryDebounced], () => {
  resetParameters()
  resetDirs()
  fetchMedias()
})

onMounted(async () => {
  await fetchMedias()
})

onMounted(async () => {
  const resp = await reqDirectories()
  if (isError(resp)) return

  for (const dir of resp.result) {
    reqDirOptions.value.push({
      label: dir.directory_alias,
      icon: dir.is_uploads ? FolderAddIcon : FolderIcon,
      directory: dir.directory_alias,
    })
  }
})
</script>
