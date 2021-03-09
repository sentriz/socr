<template>
  <table class="table-auto rounded w-full">
    <tr v-for="(directory, i) in directories" :class="{ 'bg-gray-100': i % 2 }">
      <td class="border padded">{{ directory.directory_alias }}</td>
      <td class="border padded">{{ directory.count }}</td>
    </tr>
  </table>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { isError, reqDirectories } from '../api'
import type { Directory } from '../api'

// fetch import status and about on mount
const directories = ref<Directory[]>([{directory_alias: "...", count: 0}])
onMounted(async () => {
  const resp = await reqDirectories()
  if (isError(resp)) return

  directories.value = resp.result
})
</script>
