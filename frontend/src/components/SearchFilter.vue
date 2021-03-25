<template>
  <div class="relative">
    <div class="flex border border-gray-300 bg-white rounded divide-x divide-gray-300 whitespace-nowrap">
      <div class="padded text-gray-600 bg-gray-200 rounded-l">sort by</div>
      <SearchFilterItem :label="props.selected.label" :icon="props.selected.icon" @click="toggle" />
    </div>
    <div v-if="isOpen" class="absolute z-10 right-0 py-2 mt-2 border border-gray-300 bg-white rounded text-right">
      <SearchFilterItem v-for="(item, idx) in props.items" :label="item.label" :icon="item.icon" @click="choose(idx)" />
    </div>
  </div>
</template>

<script setup lang="ts">
import SearchFilterItem from './SearchFilterItem.vue'

import { defineProps, defineEmit, ref } from 'vue'

interface Item {
  label: string
  icon: string
}

const emit = defineEmit<(e: string, v: Item) => void>()
const props = defineProps<{
  items: Item[]
  selected: Item
}>()

const isOpen = ref(false)
const toggle = () => (isOpen.value = !isOpen.value)

const choose = (index: number) => {
  emit('update:selected', props.items[index])
  toggle()
}
</script>
