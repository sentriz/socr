<template>
  <div
    ref="elm"
    class="whitespace-nowrap flex min-w-0 text-gray-700 bg-white border border-gray-300 divide-x divide-gray-300 rounded"
    :class="{ 'pointer-events-none contrast-125 text-gray-500': disabled }"
  >
    <div class="padded w-[6.5rem] bg-gray-200 rounded-l flex-shrink-0 text-right lg:text-left">
      {{ props.label }}
    </div>
    <div class="relative w-full" v-if="props.items.length">
      <search-filter-item :label="props.selected.label" :icon="props.selected.icon" @click="toggle" />

      <div v-if="isOpen" class="absolute z-10 py-2 ml-[-1px] mt-2 bg-white border border-gray-300 rounded">
        <search-filter-item
          v-for="(item, idx) in props.items"
          class="hover:bg-gray-100"
          :class="{ 'font-bold': selected === item }"
          :label="item.label"
          :icon="item.icon"
          @click="choose(idx)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import SearchFilterItem from './SearchFilterItem.vue'

import { onClickOutside } from '@vueuse/core'
import { ref } from 'vue'
import type { Component } from 'vue'

interface Item {
  label: string
  icon: Component
}

const emit = defineEmits<{ (e: 'update:selected', item: Item): void }>()
const props = defineProps<{
  label: string
  items: Item[]
  selected: Item
  disabled?: boolean
}>()

const isOpen = ref(false)
const toggle = () => (isOpen.value = !isOpen.value)
const close = () => (isOpen.value = false)

const choose = (index: number) => {
  emit('update:selected', props.items[index])
  close()
}

const elm = ref<HTMLElement>()
onClickOutside(elm, () => close())
</script>
