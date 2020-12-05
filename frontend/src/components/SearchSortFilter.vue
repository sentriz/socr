<template>
  <div
    class="flex border border-gray-400 bg-white rounded divide-x divide-gray-400 whitespace-nowrap"
  >
    <div class="padded text-gray-600 bg-gray-200 rounded-l">sort by</div>
    <div class="padded text-gray-800 w-full space-x-2 text-right" @click="toggle">
      <span class="select-none">{{ item.name }}</span>
      <i :class="item.icon"></i>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps, defineEmit } from "vue";

const emit = defineEmit<(e: 'update:modelValue', modelValue: number) => void>()
const props = defineProps<{
  modelValue: number,
  items: {
    name: string,
    icon: string,
    filter?: string[],
  }[]
}>();

const item = computed(() => props.items[props.modelValue]);
const toggle = () => {
  emit("update:modelValue", [1, 0][props.modelValue]);
};
</script>