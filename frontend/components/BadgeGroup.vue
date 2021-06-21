<template>
  <span class="whitespace-nowrap inline-flex items-center gap-2 leading-none">
    <span v-if="props.label">{{ props.label }}</span>
    <span class="inline-flex items-center space-x-2">
      <slot />
    </span>
  </span>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  label?: string
}>()
</script>
