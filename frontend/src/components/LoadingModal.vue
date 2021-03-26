<template>
  <TransitionFade>
    <div
      v-if="loading"
      class="fixed inset-0 z-10 flex items-center justify-center transition-opacity bg-gray-700 bg-opacity-75"
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
