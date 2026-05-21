<template>
  <div class="flex items-center justify-center gap-2 bg-neutral-100 p-3 text-neutral-600">
    <arrow-path-icon class="h-5 animate-spin" />
    <span>{{ props.text || 'loading' }}</span>
  </div>
</template>

<script setup lang="ts">
import { ArrowPathIcon } from '@heroicons/vue/24/outline'

const props = defineProps<{
  text?: string
}>()
</script>
