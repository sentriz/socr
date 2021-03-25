<template>
  <div v-if="screenshot" class="z-20 fixed inset-y-0 right-0 max-w-lg w-full p-6 bg-gray-100 overflow-y-auto space-y-6">
    <!-- box, header -->
    <div class="flex leading-normal">
      <router-link :to="{ name: 'search' }" class="flex-grow text-xl leading-none">
        <i class="text-gray-800 hover:text-gray-600 fas fa-times-circle"></i>
      </router-link>
      <!-- header badges -->
      <div class="flex flex-col md:flex-row gap-6 justify-end items-end">
        <BadgeGroup label="created">
          <Badge class="bg-pink-200 text-pink-900" :title="screenshot.timestamp">
            {{ relativeDateStr(screenshot.timestamp) }}
          </Badge>
        </BadgeGroup>
        <BadgeGroup v-if="screenshot.directories?.length" label="directories">
          <Badge v-for="(dir, i) in screenshot.directories" :key="i" class="bg-blue-200 text-blue-900">{{ dir }}</Badge>
        </BadgeGroup>
        <BadgeGroup>
          <Badge class="bg-indigo-200 text-indigo-900" icon="fas fa-external-link-alt">
            <a :href="screenshotRaw" target="_blank">raw</a>
          </Badge>
          <Badge class="bg-green-200 text-green-900" icon="fas fa-external-link-alt">
            <router-link :to="{ name: 'public', params: { hash: screenshot.hash } }">public</router-link>
          </Badge>
        </BadgeGroup>
      </div>
    </div>
    <!-- box -->
    <ScreenshotBackground :hash="screenshot.hash" class="box p-3">
      <ScreenshotHighlight :hash="screenshot.hash" class="mx-auto" />
    </ScreenshotBackground>
    <!-- box -->
    <div v-if="blocks.length" class="box bg-gray-200 padded font-mono text-sm">
      <p v-for="(block, i) in blocks" :key="i" :class="{ 'bg-yellow-300': highlightedBlocksIndexes.has(i) }">
        {{ block.body }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import ScreenshotHighlight from './ScreenshotHighlight.vue'
import ScreenshotBackground from './ScreenshotBackground.vue'
import BadgeGroup from './BadgeGroup.vue'
import Badge from './Badge.vue'

import { computed, defineProps, watch } from 'vue'
import relativeDate from 'relative-date'
import { urlScreenshot } from '../api/'
import useStore from '../composables/useStore'

const props = defineProps<{
  hash: string | undefined
}>()

const store = useStore()

// load the screenshot from the network if we can't find it in the store
// (can happen on page reload if we've click an image on the eg. 5th page)
watch(
  () => props.hash,
  (id) => {
    id && store.loadScreenshot(id)
  },
  { immediate: true },
)

const relativeDateStr = (stamp: string) => relativeDate(new Date(stamp))

const screenshotRaw = computed(() => `${urlScreenshot}/${props.hash}/raw`)
const screenshot = computed(() => store.getScreenshotByHash(props.hash || ''))
const blocks = computed(() => store.getBlocksByHash(props.hash || ''))
const highlightedBlocksIndexes = computed(() => {
  const hashBlocks = store.getHighlightedBlocksByHash(props.hash || '')
  return new Set(hashBlocks.map((blocks) => blocks.index))
})
</script>
