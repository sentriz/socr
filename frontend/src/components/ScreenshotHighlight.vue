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
import { urlScreenshot, Field } from "../api";
import useStore from "../composables/useStore";

const props = defineProps<{
  id: string
}>();

const store = useStore();

const screenshot = computed(() => store.getScreenshotByID(props.id));
const url = computed(() => `${urlScreenshot}/${screenshot.value.id}/raw`);

const size = computed(() => ({
  height: screenshot.value.fields[Field.SIZE_HEIGHT],
  width: screenshot.value.fields[Field.SIZE_WIDTH],
}));

const blocks = computed(() => {
  const { locations, fields } = screenshot.value;
  if (!locations?.[Field.BLOCKS_TEXT]) return [];

  const flatText = toArray(fields[Field.BLOCKS_TEXT]);
  const flatPosition = toArray(fields[Field.BLOCKS_POSITION]);
  const queriesMatches = locations[Field.BLOCKS_TEXT];
  if (!queriesMatches) return []

  const queryMatches = Object.values(queriesMatches)[0];

  return queryMatches
    .map((match) => match.array_positions)
    .flat()
    .map((i) => blockFromMatchIndexes(flatPosition, i, flatText[i]));
});

const blockFromMatchIndexes = (flatPosition: number[], i: number, text: string) => {
  const [minX, minY, maxX, maxY] = flatPosition.slice(4 * i, 4 * i + 4);
  return {
    text,
    x: minX,
    y: minY,
    width: maxX - minX,
    height: maxY - minY,
  };
};

const toArray = <T>(value: T | T[]) => (Array.isArray(value) ? value : [value]);
</script>
