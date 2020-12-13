<template>
  <div class="fixed inset-0 flex justify-center items-end p-6 pointer-events-none">
    <Transition
      enter-active-class="transform transition duration-100"
      enter-to-class="translate-y-0"
      enter-from-class="translate-y-full"
      leave-active-class="transform transition duration-100"
      leave-to-class="translate-y-full"
      leave-from-class="translate-y-0"
    >
      <Toast v-if="store.state.toast" :toast="store.state.toast"/>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import Toast from "./Toast.vue"
import useStore from "../composables/useStore"

const store = useStore()
</script>