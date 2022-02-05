<template>
  <div class="rounded-lg bg-gray-100 p-3 text-gray-700">
    after you have set socr up with your various source directories <br />
    (for example
    <span class="code">desktop</span>, <span class="code">phone</span>, <span class="code">phone recordings</span>, <span class="code">uploads</span>) <br /><br />
    then socr can import media a number of ways, including:
    <ul>
      <li>&mdash; manual upload from UI, where media is added straight to the <span class="code">uploads</span> dir</li>
      <li>&mdash; watching source directories for new media</li>
      <li>&mdash; manually scanning your selected source directories</li>
      <li class="text-gray-400">&mdash; periodic scans of the folders, coming soon</li>
    </ul>
    <br />
    as the watcher, by design, will only import new media that is sees, it can sometimes be handy to trigger a manual scan an import. so, if this is a fresh socr installation,
    click "start import" below to iterate all your socr directories and import old media 👍
  </div>
  <div class="mt-6 grid grid-cols-1 gap-3 md:grid-cols-3">
    <!-- status table -->
    <table class="col-span-full bg-white md:col-span-2">
      <colgroup>
        <col class="w-3/12" />
        <col class="w-9/12" />
      </colgroup>
      <tr>
        <td :colspan="2" class="padded border">
          <div class="space-between flex items-center justify-between py-3">
            <span v-if="status?.running && status.last_hash">
              added <span class="code">{{ status.last_hash }}</span></span
            >
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
    <div class="flex min-h-40 items-center justify-center bg-gray-100 bg-contain bg-center bg-no-repeat text-gray-500" :style="previewStyle">
      <span v-if="!url">no preview available</span>
    </div>
    <!-- errors -->
    <div class="padded col-span-full border border-red-200 bg-red-100">
      <span v-if="!status?.errors?.length" class="text-red-300">no errors yet</span>
      <ol v-if="status?.errors" v-for="error in status.errors">
        <li class="truncate text-red-900">
          {{ new Date(error.time).toLocaleTimeString() }}
          <span class="mx-3 text-red-400">|</span>
          {{ error.error }}
        </li>
      </ol>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ImportStatus } from '~/request'
import { ref, onMounted, computed, StyleValue } from 'vue'
import { newSocketAuth, urlMedia, reqStartImport, reqImportStatus, isError } from '~/request'

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

const previewStyle = computed<StyleValue>(() => {
  if (!url.value) return {}
  return { backgroundImage: `url(${url.value || ''})` }
})

// fetch import status on mount
onMounted(requestImportStatus)

// fetch import status on socket message
const socket = newSocketAuth({ want_settings: 1 })
socket.onmessage = requestImportStatus
</script>
