<template>
  <span class="px-2 space-x-1 inline-block text-white rounded">
    <i v-if="props.icon" :class="props.icon" class="text-sm" />
    <span><slot /></span>
  </span>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  icon?: string
}>()
</script>
