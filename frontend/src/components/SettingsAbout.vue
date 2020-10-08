<template>
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

import { ref, onMounted } from "vue";
import { reqAbout } from "../api";

// fetch import status and about on mount
export const about = ref({});
onMounted(async () => {
  about.value = await reqAbout();
});
</script>
