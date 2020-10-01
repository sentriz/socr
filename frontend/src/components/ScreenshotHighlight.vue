<template>
  <canvas
    ref="canvas"
    :height="scrotHeight"
    :width="scrotWidth"
    :style="{
      background: `url(${scrotURL})`,
      backgroundSize: 'cover',
    }"
    class="w-full"
  />
</template>

<script setup="props">
export default {
  props: {
    id: String,
  },
};

import { inject, ref, computed, onMounted, onUpdated } from "vue";
import { urlScreenshot, fields } from "../api";
import { zipBlocks } from "../highlighting";

export const store = inject("store");
const screenshot = computed(() => store.screenshots[props.id]);
const blocks = computed(() => zipBlocks(screenshot.value));

const highlightCanvas = (canvas) => {
  const ctx = canvas.getContext("2d");
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  for (const block of blocks.value) {
    if (!block.match) continue;

    ctx.fillStyle = "rgba(236, 201, 75, 0.75)";
    ctx.fillRect(
      block.position.minX,
      block.position.minY,
      block.position.maxX - block.position.minX,
      block.position.maxY - block.position.minY,
    );
  }
};

export const canvas = ref(null);
onMounted(() => highlightCanvas(canvas.value));
onUpdated(() => highlightCanvas(canvas.value));

export const scrotHeight = computed(() => screenshot.value.fields[fields.SIZE_HEIGHT]);
export const scrotWidth = computed(() => screenshot.value.fields[fields.SIZE_WIDTH]);
export const scrotURL = computed(() => `${urlScreenshot}/${screenshot.value.id}/raw`);
</script>
