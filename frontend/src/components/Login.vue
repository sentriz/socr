<template>
  <div class="bg-gray-200 min-h-screen flex items-center justify-center">
    <div class="w-full max-w-xs space-y-4 m-8">
      <div class="bg-white shadow-md rounded p-8 space-y-6">
        <Logo class="w-9/12 mx-auto" />
        <div>
          <label class="inp-label" for="username">username</label>
          <input
            class="inp shadow w-full"
            id="username"
            type="text"
            placeholder="mark_e_smith"
            v-model="username"
          />
        </div>
        <div>
          <label class="inp-label" for="password">password</label>
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
        <b>s</b>creenshot
        <b>ocr</b> server &mdash; Senan Kelly 2020
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import Logo from "./Logo.vue";

import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { reqAuthenticate, tokenSet } from "../api";
import type { ResponseAuthenticate } from "../api";
import useStore from "../composables/useStore"

const route = useRoute();
const router = useRouter();
const store = useStore();

const username = ref("");
const password = ref("");

const login = async () => {
  let response: ResponseAuthenticate | undefined
  try {
    response = await reqAuthenticate({
      username: username.value,
      password: password.value,
    })
  } catch (err) {
    store.setToast(err.message)
    return
  }

  if (response?.token) {
    tokenSet(response?.token);
    router.replace(route.query.redirect as string || "/");
  }
};
</script>