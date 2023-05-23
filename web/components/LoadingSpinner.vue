<template>
  <div class="flex items-center justify-center gap-2 bg-gray-100 p-3 text-gray-600">
    <refresh-icon class="h-5 animate-spin" />
    <span>{{ props.text || 'loading' }}</span>
  </div>
</template>

<script setup lang="ts">
import { RefreshIcon } from '@heroicons/vue/outline'

const props = defineProps<{
  text?: string
}>()
</script>
