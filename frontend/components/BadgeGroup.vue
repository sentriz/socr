<template>
  <span class="whitespace-nowrap">
    {{ props.label }}
    <span class="space-x-2">
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
