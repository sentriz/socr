<template>
  <div class="container mx-auto px-6 py-3">
    <Search />
  </div>
</template>

<script>
import Search from "./Search.vue";

export default {
  name: "App",
  components: {
    Search,
  },
};
</script>
