<template>
  <span class="px-2 space-x-1 inline-block text-white rounded">
    <i v-if="icon" :class="icon" class="text-sm" />
    <span><slot /></span>
  </span>
</template>

<script setup lang="ts">
import { defineProps, toRefs } from "vue";

const props = defineProps<{
  icon?: string
}>();

const { icon } = toRefs(props);
</script>
