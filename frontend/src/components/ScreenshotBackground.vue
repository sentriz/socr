<template>
  <div :style="dominantStyle">
    <slot />
  </div>
</template>

<script setup="props">
export default {
  components: {},
  props: { id: String },
};

import { computed } from "vue";
import { fields } from "../api";
import { useStore } from "../store";

const ALPHA = "88";
const store = useStore();

export const screenshot = computed(() => store.screenshotByID(props.id));
export const dominantStyle = computed(() => {
  const backgroundColor = `${screenshot.value?.fields?.[fields.DOMINANT_COLOUR]}${ALPHA}`;
  return { backgroundColor };
});
</script>
