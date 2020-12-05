<template>
  <div :style="dominantStyle">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed } from "vue";
import { fields } from "../api";
import { useStore } from "../store";

const props = defineProps<{
  id: string,
}>();

const ALPHA = "88";
const store = useStore();

const screenshot = computed(() => store.screenshotByID(props.id));
const dominantStyle = computed(() => {
  const backgroundColor = `${screenshot.value?.fields?.[fields.DOMINANT_COLOUR]}${ALPHA}`;
  return { backgroundColor };
});
</script>
