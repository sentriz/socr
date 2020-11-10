<template>
  <span class="px-2 space-x-1 inline-block text-white rounded">
    <i v-if="icon" :class="icon" class="text-sm" />
    <span><slot /></span>
  </span>
</template>

<script setup="props">
export default {
  components: {},
  props: { icon: String },
};

import { toRefs } from "vue";

export const { icon } = toRefs(props);
</script>
