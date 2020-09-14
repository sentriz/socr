<template>
  <canvas
    ref="canvas"
    :height="screenshot.fields['dimensions.height']"
    :width="screenshot.fields['dimensions.width']"
    :style="{
      background: `url(${imageURL(screenshot.id)})`,
      backgroundSize: 'cover',
    }"
  />
</template>

<script>
import { ref, onMounted } from "vue";
import { imageURL } from "../api";
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
  methods: {
    imageURL,
  },
};
</script>
