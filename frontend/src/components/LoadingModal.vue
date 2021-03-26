<template>
  <TransitionFade>
    <div
      v-if="loading"
      class="z-10 fixed inset-0 bg-gray-700 bg-opacity-75 transition-opacity flex items-center justify-center"
    >
      <LoadingSpinner :text="props.text" />
    </div>
  </TransitionFade>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'
import LoadingSpinner from './LoadingSpinner.vue'
import TransitionFade from './TransitionFade.vue'

const props = defineProps<{
  loading: boolean
  text?: string
}>()
</script>
