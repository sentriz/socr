<template>
  <div class="container">
    <input
      class="bg-white focus:outline-none focus:shadow-outline border border-gray-300 rounded-lg py-2 px-4 block w-full appearance-none leading-normal"
      type="email"
      placeholder="enter screenshot text query"
      v-model="query"
    />
    <p class="my-3 text-gray-500 text-right">
      {{ response.total_hits }} results found in {{ tookMS }}ms
    </p>
    <hr class="my-3" />
    <div id="photos">
      <Thumbnail
        v-for="screenshot in response.hits"
        :key="screenshot.id"
        :screenshot="screenshot"
        class="photo border border-gray-300 rounded-lg"
      />
    </div>
    <router-link :to="{ name: 'result', params: { id: 'wow' } }"
      >open</router-link
    >
    |
    <router-link :to="{ name: 'search' }">close</router-link>
    <router-view v-slot="{ Component, route }">
      <transition name="slide">
        <component :is="Component" v-bind="route.params"></component>
      </transition>
    </router-view>
  </div>
</template>

<script>
import { ref } from "vue";
import throttle from "lodash.debounce";

import { doSearch } from "../api";
import Thumbnail from "./Thumbnail.vue";

export default {
  data() {
    return {
      query: "",
      response: {
        hits: [],
        total_hits: 0,
        took: 0,
      },
    };
  },
  components: {
    Thumbnail,
  },
  watch: {
    query(query, _) {
      if (query) this.fetchScreenshots();
    },
  },
  computed: {
    tookMS() {
      return Math.round((this.response.took / 100000) * 100) / 100;
    },
  },
  methods: {
    fetchScreenshots: throttle(async function () {
      this.response = await doSearch({
        size: 40,
        fields: [
          "blocks.text",
          "blocks.position",
          "dimensions.height",
          "dimensions.width",
        ],
        highlight: {
          fields: ["blocks.text"],
        },
        query: {
          term: this.query,
        },
      });
    }, 200),
  },
};
</script>

<style scoped>
#photos {
  line-height: 0;
  column-count: 4;
  column-gap: 5px;
}

#photos .photo {
  width: 100%;
  height: auto;
  margin: 5px 0;
  display: flex;
  justify-content: center;
}

/* prettier-ignore */
@media (max-width: 1200px) { #photos { column-count: 3; } }
/* prettier-ignore */
@media (max-width: 1000px) { #photos { column-count: 2; } }
/* prettier-ignore */
@media (max-width: 800px)  { #photos { column-count: 1; } }

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.2s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
  transition: all 150ms ease-in 0s;
}
</style>
