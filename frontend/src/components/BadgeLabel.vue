<template>
  <span class="whitespace-nowrap">
    {{ label }}
    <span class="space-x-2">
      <slot />
    </span>
  </span>
</template>

<script setup lang="ts">
import { defineProps, toRefs } from "vue";

const props = defineProps<{
  label: string
}>();

const { label } = toRefs(props);
</script>
