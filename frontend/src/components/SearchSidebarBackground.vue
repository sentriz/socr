<template>
  <div
    v-if="screenshot"
    class="z-10 fixed inset-0 bg-gray-700 bg-opacity-75 transition-opacity pointer-events-none"
  />
</template>

<script setup lang="ts">
import { computed, defineProps } from "vue";
import useStore from "../composables/useStore";

const props = defineProps<{
  id: string | undefined,
}>();

const store = useStore();
const screenshot = computed(() => store.getScreenshotByID(props.id || ""));
</script>
