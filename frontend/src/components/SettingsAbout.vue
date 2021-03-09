<template>
  <table class="table-auto rounded w-full">
    <tr>
      <td class="border padded">version</td>
      <td class="border padded">{{ about?.version || '...' }}</td>
    </tr>
    <tr class="bg-gray-100">
      <td class="border padded">api key</td>
      <td class="border padded">{{ about?.api_key || '...' }}</td>
    </tr>
    <tr>
      <td class="border padded">socket clients</td>
      <td class="border padded">{{ about?.socket_clients || '...' }}</td>
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
