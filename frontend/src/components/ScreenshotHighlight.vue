<template>
  <div class="relative w-fit">
    <img :src="scrotURL" class="shadow-lg" />
    <canvas ref="canvas" class="absolute inset-0 w-full h-full" />
  </div>
</template>

<script setup="props">
export default {
  components: {},
  props: { id: String },
};

import { inject, ref, toRefs, computed, onMounted, onUpdated } from "vue";
import { urlScreenshot, fields } from "../api";
import { useStore } from "../store";

const store = useStore();

const screenshot = computed(() => store.screenshotByID(props.id));
const blocks = computed(() => zipBlocks(screenshot.value));
const dims = computed(() => ({
  height: screenshot.value.fields[fields.SIZE_HEIGHT],
  width: screenshot.value.fields[fields.SIZE_WIDTH],
}));

export const canvas = ref(null);

const highlightCanvas = () => {
  canvas.value.height = canvas.value.offsetHeight;
  canvas.value.width = canvas.value.offsetWidth;

  const ctx = canvas.value.getContext("2d");
  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height);

  const ratioX = canvas.value.width / dims.value.width;
  const ratioY = canvas.value.height / dims.value.height;
  for (const block of blocks.value) {
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

onMounted(() => highlightCanvas());
onUpdated(() => highlightCanvas());

export const scrotHeight = computed(() => screenshot.value.fields[fields.SIZE_HEIGHT]);
export const scrotWidth = computed(() => screenshot.value.fields[fields.SIZE_WIDTH]);
export const scrotURL = computed(() => `${urlScreenshot}/${screenshot.value.id}/raw`);

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
</script>
