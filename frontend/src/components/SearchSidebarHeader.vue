<template>
  <div class="flex leading-normal">
    <!-- left -->
    <router-link :to="{ name: 'search' }" class="flex-grow text-xl leading-none">
      <i class="hover:text-gray-600 fas fa-times-circle text-gray-800"></i>
    </router-link>
    <!-- right -->
    <div class="md:flex-row flex flex-col items-end justify-end gap-6">
      <BadgeGroup label="created">
        <Badge class="text-pink-900 bg-pink-200" :title="screenshot.timestamp">
          {{ timestampRelative }}
        </Badge>
      </BadgeGroup>
      <BadgeGroup v-if="screenshot.directories?.length" label="directories">
        <Badge v-for="(dir, i) in screenshot.directories" :key="i" class="text-blue-900 bg-blue-200">{{ dir }}</Badge>
      </BadgeGroup>
      <BadgeGroup>
        <Badge class="text-indigo-900 bg-indigo-200" icon="fas fa-external-link-alt">
          <a :href="screenshotRaw" target="_blank">raw</a>
        </Badge>
        <Badge class="text-green-900 bg-green-200" icon="fas fa-external-link-alt">
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
