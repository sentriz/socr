<template>
  <div class="transition transform bg-white text-gray-700 shadow-md border border-gray-100 rounded padded">
    {{props.toast}}
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';

const props = defineProps<{
  toast: string
}>();
</script>