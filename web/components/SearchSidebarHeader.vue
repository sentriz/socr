<template>
  <div class="flex items-center leading-normal">
    <!-- left -->
    <router-link :to="{ name: 'search' }" class="h-6 flex-grow text-xl">
      <x-icon class="h-full text-gray-800 hover:text-gray-600" />
    </router-link>
    <!-- right -->
    <div v-if="media" class="flex flex-col items-end justify-end gap-3 md:flex-row md:items-center md:gap-6">
      <badge-group label="created">
        <badge class="bg-pink-200 text-pink-900" :title="media.timestamp">
          {{ timestampRelative }}
        </badge>
      </badge-group>
      <badge-group v-if="media.directories?.length" label="directories">
        <badge v-for="(dir, i) in media.directories" :key="i" class="bg-blue-200 text-blue-900">{{ dir }}</badge>
      </badge-group>
      <badge-group>
        <a :href="mediaRaw" target="_blank">
          <badge class="bg-indigo-200 text-indigo-900">
            <external-link-icon class="h-full" />
            <span>raw</span>
          </badge>
        </a>
        <router-link :to="{ name: 'public', params: { hash: media.hash } }">
          <badge class="bg-green-200 text-green-900">
            <external-link-icon class="h-full" />
            <span>public</span>
          </badge>
        </router-link>
      </badge-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import BadgeGroup from './BadgeGroup.vue'
import Badge from './Badge.vue'

import { computed } from 'vue'
import { urlMedia } from '~/request/'
import useStore from '~/composables/useStore'
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
