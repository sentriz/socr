<template>
  <div :style="dominantStyle">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed } from "vue";
import { Field } from "../api";
import useStore from "../composables/useStore";

const props = defineProps<{
  id: string,
}>();

const ALPHA = "88";
const store = useStore();

const screenshot = computed(() => store.getScreenshotByID(props.id));
const dominantStyle = computed(() => {
  const backgroundColor = `${screenshot.value?.fields?.[Field.DOMINANT_COLOUR]}${ALPHA}`;
  return { backgroundColor };
});
</script>
