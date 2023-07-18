<template>
  <media-background ref="imgParent" v-if="media" :hash="media.hash" class="flex justify-center">
    <media-highlight ref="imgChild" :hash="media.hash" class="shadow-sm" />
  </media-background>
  <loading-spinner v-else class="bg-gray-100" text="processing image" />
</template>

<script setup lang="ts">
import MediaBackground from './MediaBackground.vue'
import MediaHighlight from './MediaHighlight.vue'
import LoadingSpinner from './LoadingSpinner.vue'
import { computed } from 'vue'
import useStore from '~/composables/useStore'

const props = defineProps<{
  hash: string
}>()

const store = useStore()
const media = computed(() => store.getMediaByHash(props.hash || ''))
</script>
