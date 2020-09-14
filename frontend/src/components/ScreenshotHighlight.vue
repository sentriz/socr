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

<script>
import { ref, onMounted } from "vue";
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

export default {
  name: "ScreenshotHighlight",
  props: {
    screenshot: Object,
  },
  setup(props) {
    const canvas = ref(null);
    const blocks = zipBlocks(props.screenshot);
    onMounted(() => {
      var ctx = canvas.value.getContext("2d");
      highlightCanvas(ctx, blocks);
    });
    return { canvas };
  },
  computed: {
    scrotHeight() {
      return this.screenshot.fields[fields.DIMENSIONS_HEIGHT];
    },
    scrotWidth() {
      return this.screenshot.fields[fields.DIMENSIONS_WIDTH];
    },
    scrotURL() {
      return `${urlImage}/${this.screenshot.id}`
    },
  },
};
</script>
