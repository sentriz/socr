<template>
  <div class="space-y-3">
    <p class="my-3 text-gray-700 text-lg font-bold">importer</p>
    <div class="flex items-center space-x-3">
      <button
        class="btn"
        :class="{ disabled: !status.finished }"
        @click="reqStartImport"
      >
        start import
      </button>
      <div class="flex-1 bg-blue-200 font-mono padded rounded">
        {{ status.new }}
      </div>
      <div v-show="!status.finished" class="bg-blue-200 padded rounded">
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
  <p class="my-3 text-gray-700 text-lg font-bold">credentials</p>
  <button class="btn">hello?</button>
</template>

<script setup>
export default {
  props: {},
};

import { ref, inject } from "vue";
export { reqStartImport, newSocketAuth } from "../api";

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
</script>
