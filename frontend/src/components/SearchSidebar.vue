<template>
  <Transition
    enter-active-class="transform transition ease-in-out duration-200"
    enter-to-class="translate-x-0"
    enter-from-class="translate-x-full"
    leave-active-class="transform transition ease-in-out duration-200"
    leave-to-class="translate-x-full"
    leave-from-class="translate-x-0"
  >
    <SearchSidebarMain :hash="props.hash" />
  </Transition>
  <Transition
    enter-active-class="ease-in-out duration-500"
    enter-to-class="opacity-100"
    enter-from-class="opacity-0"
    leave-active-class="ease-in-out duration-500"
    leave-to-class="opacity-0"
    leave-from-class="opacity-100"
  >
    <SearchSidebarBackground :hash="props.hash" />
  </Transition>
</template>

<script setup lang="ts">
import SearchSidebarMain from "./SearchSidebarMain.vue"
import SearchSidebarBackground from "./SearchSidebarBackground.vue"

import { defineProps } from "vue";

const props = defineProps<{
  hash: string | undefined,
}>();
</script>
