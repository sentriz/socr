<template>
  <Teleport to="body">
    <TransitionSlideY>
      <div v-if="store.state.toast" class="transition transform fixed bottom-6 inset-x-0 flex justify-center">
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
