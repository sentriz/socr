<template>
  <div class="bg-gray-200 h-screen flex items-center justify-center">
    <div class="w-full max-w-xs">
      <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <div class="mb-4">
          <label
            class="block text-gray-700 text-sm font-bold mb-2"
            for="username"
          >
            username
          </label>
          <input
            class="inp w-full"
            id="username"
            type="text"
            placeholder="mark_e_smith"
            v-model="username"
          />
        </div>
        <div class="mb-6">
          <label
            class="block text-gray-700 text-sm font-bold mb-2"
            for="password"
          >
            password
          </label>
          <input
            class="inp w-full"
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

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { reqAuthenticate, tokenSet } from "../api";

export const username = ref("");
export const password = ref("");

const router = useRouter();
export const login = async () => {
  try {
    const resp = await reqAuthenticate({
      username: username.value,
      password: password.value,
    });
    if (resp.token) {
      tokenSet(resp.token);
      router.push({ name: "home" });
    }
  } catch (err) {}
};
</script>
