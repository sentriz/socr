<template>
  <div class="padded flex w-full cursor-pointer items-center space-x-3">
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
