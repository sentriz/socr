<template>
  <div class="relative w-fit">
    <img :src="url" class="shadow" />
    <svg class="absolute h-full w-full"></svg>
    <!-- <canvas ref="canvas" class="absolute inset-0 w-full h-full" /> -->
    <!-- <div class="absolute inset-0 bg-red-200 bg-opacity-25">{{ id }}</div> -->
  </div>
</template>

<script setup="props">
export default {
  components: {},
  props: { id: String },
};

import {
  inject,
  ref,
  toRefs,
  computed,
  watch,
  watchEffect,
  onMounted,
  onUpdated,
} from "vue";
import { urlScreenshot, fields } from "../api";
import { useStore } from "../store";

const store = useStore();

export const screenshot = computed(() => store.screenshotByID(props.id));
export const id = computed(() => screenshot.value.id);
export const url = computed(() => `${urlScreenshot}/${screenshot.value.id}/raw`);
export const canvas = ref(null);

const highlightCanvas = (canvas, blocks, { width, height }) => {
  canvas.height = canvas.offsetHeight;
  canvas.width = canvas.offsetWidth;

  const ctx = canvas.getContext("2d");
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  const ratioX = canvas.width / width;
  const ratioY = canvas.height / height;
  for (const block of blocks) {
    if (!block.match) continue;

    ctx.fillStyle = "rgba(236, 201, 75, 0.75)";
    ctx.fillRect(
      block.position.minX * ratioX,
      block.position.minY * ratioY,
      (block.position.maxX - block.position.minX) * ratioX,
      (block.position.maxY - block.position.minY) * ratioY,
    );
  }
};

const zipBlocks = (screenshot) => {
  if (!screenshot.locations?.[fields.BLOCKS_TEXT]) return [];

  let flatText = screenshot.fields[fields.BLOCKS_TEXT];
  let flatPosition = screenshot.fields[fields.BLOCKS_POSITION];
  if (!Array.isArray(flatText)) flatText = [flatText];
  if (!Array.isArray(flatPosition)) flatPosition = [flatPosition];

  const queriesMatches = screenshot.locations[fields.BLOCKS_TEXT];
  const queryMatches = Object.values(queriesMatches)[0];
  const matchIndexes = new Set(queryMatches.map((match) => match.array_positions).flat());

  return flatText.map((block, i) => {
    const [minX, minY, maxX, maxY] = flatPosition.slice(4 * i, 4 * i + 4);
    return {
      text: block,
      position: { minX, minY, maxX, maxY },
      match: matchIndexes.has(i),
    };
  });
};

watch([canvas, screenshot], ([canvas, screenshot]) => {
  console.log(screenshot.id, Object.keys(screenshot.locations || {}).length, screenshot);
  const blocks = zipBlocks(screenshot);
  const size = {
    height: screenshot.fields[fields.SIZE_HEIGHT],
    width: screenshot.fields[fields.SIZE_WIDTH],
  };
  highlightCanvas(canvas, blocks, size);
});

// onMounted(() => highlightCanvas());
</script>
