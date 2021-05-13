<template>
  <div class="space-y-6">
    <div class="md:flex-row flex flex-col gap-2">
      <input v-model="reqQuery" class="inp w-full" type="text" placeholder="enter media text query" />
      <search-filter label="sort by" :items="reqSortOptions" v-model:selected="reqSortOption" />
      <search-filter label="directory" :items="reqFilterOptions" v-model:selected="reqFilterOption" />
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

import { watch, computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useDebounce } from '@vueuse/core'
import { isError, SortOrder, reqDirectories } from '../api'
import type { PayloadSearch } from '../api'
import useStore from '../composables/useStore'
import useInfiniteScroll from '../composables/useInfiniteScroll'
import useLoading from '../composables/useLoading'

const store = useStore()
const route = useRoute()
const { loading, load } = useLoading(store.loadMedias)

const sidebarHash = computed(() => (route.params.hash as string) || '')

const reqQuery = ref('')
const reqPageNum = ref(0)
const reqQueryDebounced = useDebounce(reqQuery, 100)
const reqPageSize = 25

type Sort = { label: string; icon: string; field: string; order: SortOrder }
const reqSortSimilarity: Sort = { label: 'similarity', icon: 'fas fa-search', field: 'similarity', order: SortOrder.Desc }
const reqSortTimeDesc: Sort = { label: 'time desc', icon: 'fas fa-chevron-down', field: 'timestamp', order: SortOrder.Desc }
const reqSortTimeAsc: Sort = { label: 'time asc', icon: 'fas fa-chevron-up', field: 'timestamp', order: SortOrder.Asc }
const reqSortOptionsDefault = [reqSortTimeDesc, reqSortTimeAsc]
const reqSortOptionsSimilarity = [reqSortSimilarity, reqSortTimeDesc, reqSortTimeAsc]
const reqSortOptions = ref(reqSortOptionsDefault)
const reqSortOption = ref(reqSortTimeDesc)

type Filter = { label: string; icon: string; directory?: string }
const reqFilterAll: Filter = { label: 'all', icon: 'fas fa-asterisk' }
const reqFilterOptions = ref([reqFilterAll])
const reqFilterOption = ref(reqFilterAll)

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
    directory: reqFilterOption.value.directory,
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

const resetFilters = () => {
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

watch([reqSortOption, reqFilterOption, reqQueryDebounced], () => {
  resetParameters()
  resetFilters()
  fetchMedias()
})

onMounted(async () => {
  await fetchMedias()
})

onMounted(async () => {
  const resp = await reqDirectories()
  if (isError(resp)) return

  for (const dir of resp.result) {
    reqFilterOptions.value.push({
      label: dir.directory_alias,
      icon: dir.is_uploads ? 'fas fa-folder-plus' : 'fas fa-folder',
      directory: dir.directory_alias,
    })
  }
})
</script>
