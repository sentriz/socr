<template>
  <h2>importer</h2>
  <div class="space-y-4">
    <button class="btn" :disabled="!status.finished" @click="reqStartImport">
      start import
    </button>
    <table>
      <tr>
        <td class="border padded">status</td>
        <td v-if="status.id" class="border padded">
          added <span class="font-mono bg-gray-300 px-2 rounded">{{ status.id }}</span>
        </td>
        <td v-else class="border padded">import not started</td>
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
      <tr v-show="errors.length">
        <td class="border padded">errors</td>
        <td class="border padded">{{ errors.length }}</td>
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

<script setup>
export default {
  props: {},
};

import { ref, inject, onMounted, computed } from "vue";
export { reqStartImport, reqAbout, newSocketAuth, urlScreenshot } from "../api";

export const progress = ref("100%");
export const url = ref("");
export const errors = ref([]);
export const status = ref({
  error: null,
  id: "",
  count_processed: 0,
  count_total: 0,
  finished: true,
});

const socket = newSocketAuth({ want_settings: 1 });
socket.onmessage = (e) => {
  try {
    status.value = JSON.parse(e.data);
  } catch (_) {
    return;
  }

  if (status.value.error) errors.value.push(status.value.error);
  if (status.value.id) url.value = `${urlScreenshot}/${status.value.id}/raw`;
  if (status.value.count_total)
    progress.value = `${Math.round(
      (100 * status.value.count_processed) / status.value.count_total,
    )}%`;
};

export const about = ref({});
onMounted(async () => {
  about.value = await reqAbout();
});
</script>
