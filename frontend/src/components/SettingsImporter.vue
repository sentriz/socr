<template>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <table class="md:col-span-2 bg-white">
      <colgroup>
        <col class="w-3/12" />
        <col class="w-9/12" />
      </colgroup>
      <tr>
        <td :colspan="2" class="border padded text-center">
          <span v-if="status.running && status.last_id">
            added
            <span class="text-sm font-mono bg-gray-300 px-2 rounded">{{ status.last_id }}</span>
          </span>
          <span v-else-if="status.running">running</span>
          <span v-else>finished</span>
        </td>
      </tr>
      <tr class="bg-gray-100">
        <td class="border padded">progress</td>
        <td class="border relative">
          <div class="z-10 absolute inset-0 bg-blue-300" :style="{ width: progress }" />
          <div class="z-20 absolute inset-0 padded text-black">{{ progress }}</div>
        </td>
      </tr>
      <tr>
        <td class="border padded">processed</td>
        <td class="border padded">{{ status.count_processed }}</td>
      </tr>
      <tr class="bg-gray-100">
        <td class="border padded">total</td>
        <td class="border padded">{{ status.count_total }}</td>
      </tr>
    </table>
    <div
      class="flex justify-center items-center text-gray-500"
      :style="{
        background: `url(${url})`,
        backgroundSize: 'contain',
        backgroundRepeat: 'no-repeat',
        backgroundPosition: 'center',
        backgroundColor: 'rgb(247, 250, 252)',
      }"
    >
      <span v-if="!url">no preview available</span>
    </div>
    <div class="md:col-span-2 bg-red-100 padded border border-red-200">
      <span v-if="!status.errors" class="text-red-300">no errors yet</span>
      <ol v-for="error in status.errors">
        <li class="text-red-900 truncate">
          {{ new Date(error.time).toLocaleTimeString() }}
          <span class="text-red-400 mx-3">|</span>
          {{ error.error }}
        </li>
      </ol>
    </div>
    <button class="btn" :disabled="status.running" @click="reqStartImport">start import</button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { newSocketAuth, urlScreenshot, reqStartImport, reqImportStatus } from "../api";
import type { ResponseImportStatus } from "../api";

const status = ref({} as ResponseImportStatus);

const requestImportStatus = async () => {
  status.value = await reqImportStatus();
};

const url = computed(() => {
  if (!status.value?.last_id) return null;
  return `${urlScreenshot}/${status.value.last_id}/raw`;
});

const progress = computed(() => {
  if (!status.value.count_total) return `0%`;
  const perc = (100 * status.value.count_processed) / status.value.count_total;
  return `${Math.round(perc)}%`;
});

// fetch import status on mount
onMounted(requestImportStatus);

// fetch import status on socket message
const socket = newSocketAuth({ want_settings: 1 });
socket.onmessage = requestImportStatus;
</script>