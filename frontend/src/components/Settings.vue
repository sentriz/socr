<template>
  <div class="container space-y-3">
    <p class="my-3 text-gray-700 text-lg font-bold">importer</p>
    <div class="flex items-center space-x-3">
      <button
        class="btn"
        :class="{ disabled: !status.finished }"
        @click="startImport"
      >
        start import
      </button>
      <div class="flex-1 bg-blue-300 font-bold text-white p-2 rounded">
        {{ status.new }}
      </div>
      <div v-show="!status.finished" class="bg-blue-300 text-white p-2 rounded">
        processed {{ status.count_processed }} / {{ status.count_total }}
      </div>
    </div>
    <div
      v-show="errors.length"
      class="flex-1 bg-red-100 border-solid border-2 border-red-200 p-2 rounded"
    >
      <p v-for="error in errors" :key="error"><b>error:</b> {{ error }}</p>
    </div>
  </div>
  <hr />
  <p class="my-3 text-gray-700 text-lg font-bold">something else</p>
  <button class="btn">hello?</button>
</template>

<script>
import { reqStartImport } from "../api";

export default {
  name: "Settings",
  inject: ["socket"],
  data() {
    return {
      errors: [],
      status: {
        error: null,
        new: "import not started",
        count_processed: 0,
        count_total: 0,
        finished: true,
      },
    };
  },
  created() {
    this.socket.onmessage = (e) => {
      this.status = JSON.parse(e.data);
      if (this.status.error) {
        this.errors.push(this.status.error);
      }
    };
  },
  methods: {
    startImport() {
      reqStartImport();
    },
  },
};
</script>
