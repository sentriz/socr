<template>
    <span class="select-none">{{ props.name }}</span>
    <i :class="props.icon"></i>
</template>

<script setup lang="ts">
import { defineProps } from "vue";

const props = defineProps<{
  id: string,
  name: string,
  icon: string
}>();
</script>