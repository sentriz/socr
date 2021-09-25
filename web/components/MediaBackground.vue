<template>
  <div :style="dominantStyle">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import useStore from '~/composables/useStore'

const props = defineProps<{
  hash?: string
}>()

const ALPHA = '88'
const store = useStore()

const media = computed(() => store.getMediaByHash(props.hash || ''))
const dominantStyle = computed(() => {
  const backgroundColor = `${media.value?.dominant_colour}${ALPHA}`
  return { backgroundColor }
})
</script>
