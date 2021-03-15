<template>
  <div class="flex border border-gray-300 bg-white rounded divide-x divide-gray-300 whitespace-nowrap">
    <div class="padded text-gray-600 bg-gray-200 rounded-l">sort by</div>
    <div class="padded text-gray-800 w-full space-x-2 text-right" @click="toggle">
      <span class="select-none">{{ props?.label }}</span>
      <i :class="icons[props?.order]"></i>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmit } from 'vue'

// TODO: improve the defineEmit types
// something like
// const emit = defineEmit<
//   ((e: 'update:field', v: string) => void) |
//   ((e: 'update:order', v: Order) => void)
// >()
// see https://github.com/vuejs/vue-next/issues/2874
// see https://github.com/vuejs/vue-next/pull/2878

const emit = defineEmit<(e: string, v: string) => void>()
const props = defineProps<{
  label: string
  field: string
  order: Order
}>()

enum Order {
  Asc = 'asc',
  Desc = 'desc',
}

const icons: { [key in Order]: string } = {
  [Order.Asc]: 'fas fa-chevron-up',
  [Order.Desc]: 'fas fa-chevron-down',
}

const toggle = () => {
  emit('update:field', props.field)
  emit('update:order', props.order === Order.Asc ? Order.Desc : Order.Asc)
}
</script>
