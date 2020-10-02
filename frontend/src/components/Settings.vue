<template>
  <h2>importer</h2>
  <div class="space-y-3">
    <div class="flex items-center space-x-3">
      <button class="btn" :disabled="!status.finished" @click="reqStartImport">
        start import
      </button>
      <div class="flex-1 bg-blue-100 font-mono padded rounded">
        {{ status.new }}
      </div>
      <div v-show="!status.finished" class="bg-blue-100 padded rounded">
        processed {{ status.count_processed }} / {{ status.count_total }}
      </div>
    </div>
    <div
      v-show="errors.length"
      class="flex-1 bg-red-100 border-solid border-2 border-red-200 padded rounded"
    >
      <p v-for="error in errors" :key="error"><b>error:</b> {{ error }}</p>
    </div>
  </div>
  <hr />
  <h2>about</h2>
  <table class="table-auto">
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

import { ref, inject, onMounted } from "vue";
export { reqStartImport, reqAbout, newSocketAuth } from "../api";

export const errors = ref([]);
export const status = ref({
  error: null,
  new: "import not started",
  count_processed: 0,
  count_total: 0,
  finished: true,
});

const socket = newSocketAuth({ want_settings: 1 });
socket.onmessage = (e) => {
  status.value = JSON.parse(e.data);
  if (status.value.error) {
    errors.value.push(status.value.error);
  }
};

export const about = ref({});
onMounted(async () => {
  about.value = await reqAbout();
});
</script>
