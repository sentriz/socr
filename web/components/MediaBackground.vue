<template>
  <div :style="dominantStyle">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed, CSSProperties } from 'vue'
import useStore from '~/composables/useStore'

const props = defineProps<{
  hash?: string
}>()

const store = useStore()
const media = computed(() => store.getMediaByHash(props.hash || ''))

const dominantStyle = computed((): CSSProperties => {
  const ALPHA = 0.25
  return {
    backgroundColor: `${media.value?.dominant_colour}${hex(ALPHA)}`,
  }
})

const hex = (v: number /* 0-1 */): string => (v * 256).toString(16)
</script>
