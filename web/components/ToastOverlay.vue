<template>
  <transition-slide-y>
    <div v-if="store.state.toast" class="absolute inset-x-0 bottom-6 flex justify-center transition">
      <toast :toast="store.state.toast" />
    </div>
  </transition-slide-y>
</template>

<script setup lang="ts">
import TransitionSlideY from './TransitionSlideY.vue'

import Toast from './Toast.vue'
import useStore from '~/composables/useStore'

const store = useStore()
</script>
