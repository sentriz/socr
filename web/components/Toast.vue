<template>
  <div class="padded inline-block text-gray-700 bg-white border border-gray-100 rounded shadow-md">
    {{ props.toast }}
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  toast: string
}>()
</script>
