<template>
  <div class="p-3 text-center text-gray-600 bg-gray-200 rounded">
    <i class="animate-spin fas fa-circle-notch mr-1"></i> {{ props.text || 'loading' }}
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  text?: string
}>()
</script>
