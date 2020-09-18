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
import { ref, toRefs, onMounted, computed } from "vue";
import { urlImage, fields as afields } from "../api";
import { zipBlocks } from "../highlighting";

export default {
  props: {
    screenshot: Object,
  },
};

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
onMounted(() => {
  const ctx = canvas.value.getContext("2d");
  const blocks = zipBlocks(props.screenshot);
  highlightCanvas(ctx, blocks);
});

const { fields, id } = toRefs(props.screenshot);
export const scrotHeight = computed(() => fields.value[afields.SIZE_HEIGHT]);
export const scrotWidth = computed(() => fields.value[afields.SIZE_WIDTH]);
export const scrotURL = computed(() => `${urlImage}/${id.value}`);
</script>
