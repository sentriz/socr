<template>
  <transition-fade>
    <div
      v-if="loading"
      class="fixed inset-0 z-10 flex items-center justify-center transition-opacity bg-gray-700 bg-opacity-75"
    >
      <loading-spinner :text="props.text" />
    </div>
  </transition-fade>
</template>

<script setup lang="ts">
import LoadingSpinner from './LoadingSpinner.vue'
import TransitionFade from './TransitionFade.vue'

const props = defineProps<{
  loading: boolean
  text?: string
}>()
</script>
