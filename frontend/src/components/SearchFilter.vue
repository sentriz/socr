<template>
  <div ref="elm" class="flex border border-gray-300 bg-white rounded divide-x divide-gray-300 whitespace-nowrap">
    <div class="padded text-gray-600 bg-gray-200 rounded-l">{{ props.label }}</div>
    <div class="relative" v-if="props.items.length">
      <SearchFilterItem :label="props.selected.label" :icon="props.selected.icon" @click="toggle" />
      <div v-if="isOpen" class="absolute z-10 py-2 mt-2 border border-gray-300 bg-white rounded">
        <SearchFilterItem
          class="hover:bg-gray-100"
          v-for="(item, idx) in props.items"
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
import { defineProps, defineEmit, ref } from 'vue'

interface Item {
  label: string
  icon: string
}

const emit = defineEmit<(e: string, v: Item) => void>()
const props = defineProps<{
  label: string
  items: Item[]
  selected: Item
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
