<template>
  <div
    class="border border-gray-400 bg-white rounded flex divide-x divide-gray-400 whitespace-no-wrap"
  >
    <div class="flex items-center padded bg-gray-200 text-gray-600">sort by</div>
    <div
      @click="toggle"
      class="flex items-center padded justify-end w-full space-x-2 text-gray-800"
    >
      <span class="select-none">{{ item.name }}</span>
      <i :class="`fas ${item.icon}`"></i>
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
