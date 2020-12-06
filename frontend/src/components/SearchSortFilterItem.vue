<template>
    <span class="select-none">{{ name }}</span>
    <i :class="icon"></i>
</template>

<script setup lang="ts">
import { defineProps, toRefs } from "vue";

const props = defineProps<{
  id: string,
  name: string,
  icon: string
}>();

const { name, icon } = toRefs(props)
</script>