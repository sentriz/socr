<template>
  <div class="relative w-fit">
    <img :src="url" />
    <svg
      :viewBox="`0 0 ${size.width} ${size.height}`"
      class="absolute inset-0 fill-current text-yellow-500 text-opacity-50"
    >
      <rect
        v-for="(block, i) in blocks"
        :key="i"
        :x="block.x"
        :y="block.y"
        :height="block.height"
        :width="block.width"
      />
    </svg>
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed } from "vue";
import { urlScreenshot, fields as apifields } from "../api";
import { useStore } from "../store";

const props = defineProps<{
  id: string
}>();

const store = useStore();

const screenshot = computed(() => store.screenshotByID(props.id));
const id = computed(() => screenshot.value.id);
const url = computed(() => `${urlScreenshot}/${screenshot.value.id}/raw`);

const size = computed(() => ({
  height: screenshot.value.fields[apifields.SIZE_HEIGHT],
  width: screenshot.value.fields[apifields.SIZE_WIDTH],
}));

const blocks = computed(() => {
  const { locations, fields } = screenshot.value;
  if (!locations?.[apifields.BLOCKS_TEXT]) return [];

  const flatText = toArray(fields[apifields.BLOCKS_TEXT]);
  const flatPosition = toArray(fields[apifields.BLOCKS_POSITION]);

  const queriesMatches = locations[apifields.BLOCKS_TEXT];
  const queryMatches = Object.values(queriesMatches)[0];

  return queryMatches
    .map((match) => match.array_positions)
    .flat()
    .map((i) => blockFromMatchIndexes(flatPosition, i, flatText[i]));
});

const blockFromMatchIndexes = (flatPosition, i, text) => {
  const [minX, minY, maxX, maxY] = flatPosition.slice(4 * i, 4 * i + 4);
  return {
    text,
    x: minX,
    y: minY,
    width: maxX - minX,
    height: maxY - minY,
  };
};

const toArray = (value) => (Array.isArray(value) ? value : [value]);
</script>
