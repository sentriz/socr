<template>
  <media-background v-if="media" :hash="media.hash" class="flex justify-center shadow-inner">
    <media-highlight :hash="media.hash" class="shadow-sm" v-bind="$attrs" />
  </media-background>
  <loading-spinner v-else class="bg-gray-100" text="processing image" />
</template>

<script setup lang="ts">
import MediaBackground from './MediaBackground.vue'
import MediaHighlight from './MediaHighlight.vue'
import LoadingSpinner from './LoadingSpinner.vue'
import { computed } from 'vue'
import useStore from '~/composables/useStore'

defineOptions({
  inheritAttrs: false,
})

const props = defineProps<{
  hash: string
}>()

const store = useStore()
const media = computed(() => store.getMediaByHash(props.hash || ''))
</script>
