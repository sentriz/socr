<template>
  <div class="flex items-center leading-normal">
    <!-- left -->
    <router-link :to="{ name: 'search' }" class="flex-grow h-6 text-xl">
      <x-icon class="hover:text-gray-600 h-full text-gray-800" />
    </router-link>
    <!-- right -->
    <div v-if="media" class="md:flex-row md:items-center md:gap-6 flex flex-col items-end justify-end gap-3">
      <badge-group label="created">
        <badge class="text-pink-900 bg-pink-200" :title="media.timestamp">
          {{ timestampRelative }}
        </badge>
      </badge-group>
      <badge-group v-if="media.directories?.length" label="directories">
        <badge v-for="(dir, i) in media.directories" :key="i" class="text-blue-900 bg-blue-200">{{ dir }}</badge>
      </badge-group>
      <badge-group>
        <badge class="text-indigo-900 bg-indigo-200">
          <external-link-icon class="h-full" />
          <a :href="mediaRaw" target="_blank">raw</a>
        </badge>
        <badge class="text-green-900 bg-green-200">
          <external-link-icon class="h-full" />
          <router-link :to="{ name: 'public', params: { hash: media.hash } }">public</router-link>
        </badge>
      </badge-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import BadgeGroup from './BadgeGroup.vue'
import Badge from './Badge.vue'

import { computed } from 'vue'
import { urlMedia } from '../api/'
import useStore from '../composables/useStore'
import { useTimeAgo } from '@vueuse/core'
import { XIcon, ExternalLinkIcon } from '@heroicons/vue/outline'

const props = defineProps<{
  hash: string
}>()

const store = useStore()

const mediaRaw = computed(() => `${urlMedia}/${props.hash}/raw`)
const media = computed(() => store.getMediaByHash(props.hash || ''))

const timestamp = computed(() => media.value?.timestamp)
const timestampRelative = useTimeAgo(timestamp)
</script>
