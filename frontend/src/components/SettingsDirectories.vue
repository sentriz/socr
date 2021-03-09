<template>
  <table class="table-auto rounded">
    <tr v-for="(directory, i) in directories">
      <td class="border padded">{{ directory.directory_alias || '...' }}</td>
      <td class="border padded">{{ directory.count || '...' }}</td>
    </tr>
  </table>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { isError, reqDirectories } from '../api'
import type { Directory } from '../api'

// fetch import status and about on mount
const directories = ref<Directory[] | undefined>()
onMounted(async () => {
  const resp = await reqDirectories()
  if (isError(resp)) return

  directories.value = resp.result
})
</script>
