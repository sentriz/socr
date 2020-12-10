<template>
  <div class="flex border border-gray-300 bg-white rounded divide-x divide-gray-300 whitespace-nowrap">
    <div class="padded text-gray-600 bg-gray-200 rounded-l">sort by</div>
    <div class="padded text-gray-800 w-full space-x-2 text-right" @click="toggle">
      <span class="select-none">{{ item.status }}</span>
      <i :class="item.icon"></i>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps, defineEmit } from "vue";

const emit = defineEmit<(e: 'update:modelValue', modelValue: string) => void>()
const props = defineProps<{
  modelValue: string,
  values: { [key: string]: { icon: string, status: string } }
}>();

const item = computed(() => props.values[props.modelValue]);
const toggle = () => {
  const items = Object.keys(props.values)
  const curr = items.indexOf(props.modelValue)
  emit("update:modelValue", items[(curr + 1) % items.length])
}
</script>