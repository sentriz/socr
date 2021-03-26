<template>
  <table class="w-full rounded table-auto">
    <tr v-for="(directory, i) in directories" :class="{ 'bg-gray-100': i % 2 }">
      <td class="padded border">{{ directory.directory_alias }}</td>
      <td class="padded border">{{ directory.count }}</td>
    </tr>
    <tr v-if="!directories.length">
      <td class="padded border">none yet</td>
    </tr>
  </table>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { isError, reqDirectories } from '../api'
import type { Directory } from '../api'

// fetch import status and about on mount
const directories = ref<Directory[]>([])
onMounted(async () => {
  const resp = await reqDirectories()
  if (isError(resp)) return

  directories.value = resp.result
})
</script>
