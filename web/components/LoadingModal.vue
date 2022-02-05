<template>
  <transition-fade>
    <div v-if="loading" class="absolute inset-0 z-10 flex items-center justify-center bg-gray-700 bg-opacity-75 transition-opacity">
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
