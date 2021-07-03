<template>
  <div class="padded flex items-center w-full space-x-3 text-gray-800 cursor-pointer">
    <component :is="props.icon" class="h-5" />
    <span class="select-none">{{ props.label }}</span>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

const props = defineProps<{
  label: string
  icon: Component
}>()
</script>
