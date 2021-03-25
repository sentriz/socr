<template>
  <div class="flex leading-normal">
    <!-- left -->
    <router-link :to="{ name: 'search' }" class="flex-grow text-xl leading-none">
      <i class="text-gray-800 hover:text-gray-600 fas fa-times-circle"></i>
    </router-link>
    <!-- right -->
    <div class="flex flex-col md:flex-row gap-6 justify-end items-end">
      <BadgeGroup label="created">
        <Badge class="bg-pink-200 text-pink-900" :title="screenshot.timestamp">
          {{ timestampRelative }}
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
</template>

<script setup lang="ts">
import BadgeGroup from './BadgeGroup.vue'
import Badge from './Badge.vue'

import { computed, defineProps } from 'vue'
import { urlScreenshot } from '../api/'
import useStore from '../composables/useStore'
import { useTimeAgo } from '@vueuse/core'

const props = defineProps<{
  hash: string
}>()

const store = useStore()

const screenshotRaw = computed(() => `${urlScreenshot}/${props.hash}/raw`)
const screenshot = computed(() => store.getScreenshotByHash(props.hash || ''))

const timestamp = computed(() => screenshot.value?.timestamp)
const timestampRelative = useTimeAgo(timestamp)
</script>
