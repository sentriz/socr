<template>
  <div class="padded text-gray-800 w-full space-x-2 text-right">
    <span class="select-none">{{ props.label }}</span>
    <i :class="props.icon"></i>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  label: string
  icon: string
}>()
</script>
