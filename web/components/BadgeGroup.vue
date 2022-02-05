<template>
  <span class="inline-flex items-center gap-2 whitespace-nowrap leading-none">
    <span v-if="props.label">{{ props.label }}</span>
    <span class="inline-flex items-center space-x-2">
      <slot />
    </span>
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{
  label?: string
}>()
</script>
