<template>
  <table class="w-full rounded table-auto">
    <tr v-if="about" v-for="(value, key, i) in about" :class="{ 'bg-gray-100': i % 2 }">
      <td class="padded border">{{ key }}</td>
      <td class="padded border">{{ value }}</td>
    </tr>
  </table>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { isError, reqAbout } from '../api'
import type { About } from '../api'

// fetch import status and about on mount
const about = ref<About | undefined>()
onMounted(async () => {
  const resp = await reqAbout()
  if (isError(resp)) return

  about.value = resp.result
})
</script>
