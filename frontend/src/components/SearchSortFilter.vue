<template>
  <div class="flex border border-gray-400 bg-white rounded divide-x divide-gray-400 whitespace-no-wrap">
    <div class="padded text-gray-600 bg-gray-200 rounded-l">sort by</div>
    <div class="padded text-gray-800 w-full space-x-2 text-right" @click="toggle">
      <span class="select-none">{{ item.name }}</span>
      <i :class="item.icon"></i>
    </div>
  </div>
</template>

<script setup="props, { emit }">
export default {
  components: {},
  props: {
    items: Array,
    modelValue: Number,
  },
};

import { computed } from "vue";

export const item = computed(() => props.items[props.modelValue]);
export const toggle = () => {
  emit("update:modelValue", (props.modelValue === 0) | 0);
};
</script>
