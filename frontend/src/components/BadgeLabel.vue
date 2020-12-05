<template>
  <span class="whitespace-nowrap">
    {{ label }}
    <span class="space-x-2">
      <slot />
    </span>
  </span>
</template>

<script setup="props">
export default {
  components: {},
  props: { label: String },
};

import { toRefs } from "vue";

export const { label } = toRefs(props);
</script>
