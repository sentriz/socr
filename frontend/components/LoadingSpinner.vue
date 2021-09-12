<template>
  <div class="flex items-center justify-center gap-2 p-3 text-gray-600 bg-gray-200 rounded">
    <refresh-icon class="animate-spin h-5" />
    <span>{{ props.text || 'loading' }}</span>
  </div>
</template>

<script setup lang="ts">
import { RefreshIcon } from '@heroicons/vue/outline'

const props = defineProps<{
  text?: string
}>()
</script>
