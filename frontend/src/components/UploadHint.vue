<template>
  <input class="hidden" ref="elm" type="file" />
  <div class="fixed bottom-0 opacity-75 hover:opacity-100 right-0 padded bg-gray-200 rounded-tl">
    <div class="flex items-center gap-2 text-gray-600 hover:text-gray-900 hover:cursor-pointer" @click="click">
      <i class="fas fa-upload text-sm"></i>
      <span>select file to upload or just paste</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const elm = ref<HTMLElement>()
const click = () => elm.value?.click()
</script>
