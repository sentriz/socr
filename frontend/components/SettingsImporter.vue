<template>
  <div class="md:grid-cols-3 grid grid-cols-1 gap-3">
    <!-- status table -->
    <table class="col-span-full md:col-span-2 bg-white">
      <colgroup>
        <col class="w-3/12" />
        <col class="w-9/12" />
      </colgroup>
      <tr>
        <td :colspan="2" class="padded border">
          <div class="space-between flex items-center justify-between py-3">
            <span v-if="status?.running && status.last_hash">
              added
              <span class="px-2 font-mono text-sm bg-gray-300 rounded">{{ status.last_hash }}</span>
            </span>
            <span v-else-if="status?.running">running</span>
            <span v-else>finished</span>
            <button class="btn" :disabled="status?.running" @click="reqStartImport">start import</button>
          </div>
        </td>
      </tr>
      <tr class="bg-gray-100">
        <td class="padded border">progress</td>
        <td class="relative border">
          <div class="absolute inset-0 z-10 bg-blue-300" :style="{ width: progress }" />
          <div class="padded absolute inset-0 z-20 text-black">{{ progress }}</div>
        </td>
      </tr>
      <tr>
        <td class="padded border">processed</td>
        <td class="padded border">{{ status?.count_processed || 0 }}</td>
      </tr>
      <tr class="bg-gray-100">
        <td class="padded border">total</td>
        <td class="padded border">{{ status?.count_total || 0 }}</td>
      </tr>
    </table>
    <!-- preview window -->
    <div
      class="min-h-40 flex items-center justify-center text-gray-500 bg-gray-100 bg-center bg-no-repeat bg-contain"
      :style="{ backgroundImage: `url(${url})` }"
    >
      <span v-if="!url">no preview available</span>
    </div>
    <!-- errors -->
    <div class="col-span-full padded bg-red-100 border border-red-200">
      <span v-if="!status?.errors?.length" class="text-red-300">no errors yet</span>
      <ol v-if="status?.errors" v-for="error in status.errors">
        <li class="text-red-900 truncate">
          {{ new Date(error.time).toLocaleTimeString() }}
          <span class="mx-3 text-red-400">|</span>
          {{ error.error }}
        </li>
      </ol>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { newSocketAuth, urlMedia, reqStartImport, reqImportStatus, isError } from '../api'
import type { ImportStatus } from '../api'

const status = ref<ImportStatus | undefined>()

const requestImportStatus = async () => {
  const resp = await reqImportStatus()
  if (isError(resp)) return

  status.value = resp.result
}

const url = computed(() => {
  if (!status.value?.last_hash) return null
  return `${urlMedia}/${status.value.last_hash}/raw`
})

const progress = computed(() => {
  if (!status.value?.count_total) return `0%`
  const perc = (100 * status.value.count_processed) / status.value.count_total
  return `${Math.round(perc)}%`
})

// fetch import status on mount
onMounted(requestImportStatus)

// fetch import status on socket message
const socket = newSocketAuth({ want_settings: 1 })
socket.onmessage = requestImportStatus
</script>
