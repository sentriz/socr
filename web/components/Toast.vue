<template>
  <div class="padded inline-block rounded border border-gray-100 bg-white text-gray-700 shadow-md">
    {{ props.toast }}
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  toast: string
}>()
</script>
