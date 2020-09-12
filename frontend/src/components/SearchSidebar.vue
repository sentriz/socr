<template>
  <div
    class="sidebar border-l-4 p-6 sidebar bg-blue-200 fixed h-full top-0 right-0"
  >
    <img class="w-full" :src="imageURL(screenshot.id)" />
  </div>
</template>

<script>
import { imageURL } from "../api/";

export default {
  name: "SearchSidebar",
  props: {
    id: String,
    results: Array,
  },
  computed: {
    // TODO: not pass all results to this component
    // perhaps use vuex
    screenshot() {
      return this.results.find((result) => result.id == this.id);
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
