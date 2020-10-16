<template>
  <div class="bg-gray-200 h-screen flex items-center justify-center">
    <div class="w-full max-w-xs space-y-4 m-8">
      <div class="bg-white shadow-md rounded p-8 space-y-6">
        <Logo class="w-9/12 mx-auto" />
        <div>
          <label class="inp-label" for="username"> username </label>
          <input
            class="inp shadow w-full"
            id="username"
            type="text"
            placeholder="mark_e_smith"
            v-model="username"
          />
        </div>
        <div>
          <label class="inp-label" for="password"> password </label>
          <input
            class="inp shadow w-full"
            id="password"
            type="password"
            placeholder="*******"
            v-model="password"
          />
        </div>
        <button class="btn w-full" type="button" @click="login">sign in</button>
      </div>
      <p class="text-center text-gray-500 text-xs">
        <b>s</b>creenshot <b>ocr</b> server &mdash; Senan Kelly 2020
      </p>
    </div>
  </div>
</template>

<script setup="props">
import Logo from "./Logo.vue";
export default {
  components: { Logo },
  props: {},
};

import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { reqAuthenticate, tokenSet } from "../api";

const route = useRoute();
const router = useRouter();

export const username = ref("");
export const password = ref("");
export const login = async () => {
  try {
    const resp = await reqAuthenticate({
      username: username.value,
      password: password.value,
    });
    if (resp.token) {
      tokenSet(resp.token);
      router.replace(route.query.redirect || "/");
    }
  } catch (err) {}
};
</script>
