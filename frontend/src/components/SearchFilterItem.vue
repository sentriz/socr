<template>
  <div class="padded text-gray-800 w-full space-x-3 cursor-pointer">
    <i :class="props.icon"></i>
    <span class="select-none">{{ props.label }}</span>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  label: string
  icon: string
}>()
</script>
