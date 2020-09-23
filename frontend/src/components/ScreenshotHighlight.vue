<template>
  <canvas
    ref="canvas"
    :height="scrotHeight"
    :width="scrotWidth"
    :style="{
      background: `url(${scrotURL})`,
      backgroundSize: 'cover',
    }"
  />
</template>

<script setup="props">
export default {
  props: {
    screenshot: Object,
  },
};

import { ref, toRefs, onMounted, computed, watch } from "vue";
import { urlImage, fields } from "../api";
import { zipBlocks } from "../highlighting";

const highlightCanvas = (ctx, blocks) => {
  for (const block of blocks) {
    if (!block.match) continue;

    ctx.fillStyle = "rgba(236, 201, 75, 0.75)";
    ctx.fillRect(
      block.position.minX,
      block.position.minY,
      block.position.maxX - block.position.minX,
      block.position.maxY - block.position.minY
    );
  }
};

export const canvas = ref(null);
const blocks = zipBlocks(props.screenshot);
onMounted(() => {
  const ctx = canvas.value.getContext("2d");
  highlightCanvas(ctx, blocks);
});

export const scrotHeight = computed(
  () => props.screenshot.fields[fields.SIZE_HEIGHT]
);
export const scrotWidth = computed(
  () => props.screenshot.fields[fields.SIZE_WIDTH]
);
export const scrotURL = computed(() => `${urlImage}/${props.screenshot.id}`);
</script>
