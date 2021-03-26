<template>
  <Teleport to="body">
    <TransitionSlideY>
      <div v-if="store.state.toast" class="bottom-6 fixed inset-x-0 flex justify-center transition transform">
        <Toast :toast="store.state.toast" />
      </div>
    </TransitionSlideY>
  </Teleport>
</template>

<script setup lang="ts">
import TransitionSlideY from './TransitionSlideY.vue'

import Toast from './Toast.vue'
import useStore from '../composables/useStore'

const store = useStore()
</script>
