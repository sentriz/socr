<template>
  <div :style="dominantStyle">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed } from "vue";
import useStore from "../composables/useStore";

const props = defineProps<{
  hash: string,
}>();

const ALPHA = "88";
const store = useStore();

const screenshot = computed(() => store.getScreenshotByHash(props.hash));
const dominantStyle = computed(() => {
  const backgroundColor = `${screenshot.value?.dominant_colour}${ALPHA}`;
  return { backgroundColor };
});
</script>
