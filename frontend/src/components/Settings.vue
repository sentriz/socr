<template>
  <h2>importer</h2>
  <div class="space-y-4">
    <button class="btn" :disabled="status.running" @click="reqStartImport">
      start import
    </button>
    <table>
      <tr>
        <td class="border padded">status</td>
        <td class="border padded">
          <span v-if="status.running && status.last_id">
            added
            <span class="font-mono bg-gray-300 px-2 rounded">{{ status.last_id }}</span>
          </span>
          <span v-else-if="status.running">running</span>
          <span v-else>finished</span>
        </td>
        <td
          v-show="url"
          class="hidden lg:table-cell w-64"
          rowspan="0"
          :style="{
            background: `url(${url})`,
            backgroundSize: 'contain',
            backgroundRepeat: 'no-repeat',
            backgroundPosition: 'left',
          }"
        ></td>
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
      <tr v-show="status.errors?.length">
        <td class="border padded">errors</td>
        <td class="border padded">{{ status.errors?.length }}</td>
      </tr>
    </table>
  </div>
  <hr />
  <h2>about</h2>
  <table class="table-auto rounded">
    <tr>
      <td class="border padded">version</td>
      <td class="border padded">{{ about.version || "..." }}</td>
    </tr>
    <tr class="bg-gray-100">
      <td class="border padded">screenshots indexed</td>
      <td class="border padded">{{ about.screenshots_indexed || "..." }}</td>
    </tr>
    <tr>
      <td class="border padded">api key</td>
      <td class="border padded">{{ about.api_key || "..." }}</td>
    </tr>
    <tr class="bg-gray-100">
      <td class="border padded">import path</td>
      <td class="border padded">{{ about.import_path || "..." }}</td>
    </tr>
    <tr>
      <td class="border padded">screenshots path</td>
      <td class="border padded">{{ about.screenshots_path || "..." }}</td>
    </tr>
  </table>
</template>

<script setup="props">
export default {
  components: {},
  props: {},
};

import { ref, inject, onMounted, computed } from "vue";
import {  newSocketAuth, urlScreenshot } from "../api";
export { reqStartImport, reqImportStatus, reqAbout } from "../api"

export const status = ref({
  running: false,
  errors: [],
  last_id: "",
  count_processed: 0,
  count_total: 0,
});

export const url = computed(() => {
  if (!status.value.last_id) return null;
  return `${urlScreenshot}/${status.value.last_id}/raw`;
});

export const progress = computed(() => {
  if (!status.value.count_total) return `0%`;
  const perc = (100 * status.value.count_processed) / status.value.count_total;
  return `${Math.round(perc)}%`;
});

// fetch import status and about on mount
export const about = ref({});
onMounted(async () => {
  status.value = await reqImportStatus();
  about.value = await reqAbout();
});

// fetch import status on socket message
const socket = newSocketAuth({ want_settings: 1 });
socket.onmessage = async () => {
  status.value = await reqImportStatus();
};
</script>
