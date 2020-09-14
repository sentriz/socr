<template>
  <div
    class="sidebar border-l-4 p-6 sidebar bg-white fixed h-full top-0 right-0"
  >
    <div class="mx-auto">
      <div class="bg-black shadow font-mono text-sm">
        <ScreenshotHighlight class="mx-auto" :screenshot="screenshot" />
      </div>
      <hr class="my-6" />
      <div class="bg-gray-300 p-3 shadow font-mono text-sm">
        <p v-for="(line, i) in text" :key="i">
          {{ line }}
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { imageURL } from "../api/";
import ScreenshotHighlight from "./ScreenshotHighlight.vue";

export default {
  name: "SearchSidebar",
  props: {
    id: String,
    results: Array,
  },
  components: {
    ScreenshotHighlight,
  },
  computed: {
    // TODO: not pass all results to this component
    // perhaps use vuex
    screenshot() {
      return this.results.find((result) => result.id == this.id);
    },
    text() {
      return this.screenshot.fields["blocks.text"];
    },
  },
  methods: {
    imageURL,
  },
};
</script>

<style scoped>
.sidebar {
  width: 75vw;
}
</style>
