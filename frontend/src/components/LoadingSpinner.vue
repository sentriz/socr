<template>
  <div class="bg-gray-200 text-gray-600 text-center rounded p-3">
    <i class="animate-spin fas fa-circle-notch mr-1"></i> {{ props.text || 'loading' }}
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  text?: string
}>()
</script>
