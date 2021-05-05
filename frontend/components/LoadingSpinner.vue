<template>
  <div class="flex items-center justify-center gap-2 p-3 text-gray-600 bg-gray-200 rounded">
    <i class="animate-spin fas fa-circle-notch" />
    <span>{{ props.text || 'loading' }}</span>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  text?: string
}>()
</script>
