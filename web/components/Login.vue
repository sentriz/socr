<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-200">
    <div class="w-full max-w-xs m-8 space-y-4">
      <div class="p-8 space-y-6 bg-white rounded shadow-md">
        <logo class="w-9/12 mx-auto" />
        <div>
          <label class="inp-label" for="username">username</label>
          <input class="inp w-full shadow" hash="username" type="text" placeholder="mark_e_smith" v-model="username" />
        </div>
        <div>
          <label class="inp-label" for="password">password</label>
          <input class="inp w-full shadow" hash="password" type="password" placeholder="*******" v-model="password" />
        </div>
        <button class="btn w-full" type="button" @click="login">sign in</button>
      </div>
      <p class="text-xs text-center text-gray-500"><b>s</b>creenshot <b>ocr</b> server &mdash; Senan Kelly 2020</p>
    </div>
  </div>
  <teleport to="#overlays">
    <toast-overlay />
  </teleport>
</template>

<script setup lang="ts">
import Logo from './Logo.vue'
import ToastOverlay from './ToastOverlay.vue'

import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { isError, reqAuthenticate, tokenSet } from '~/request'
import useStore from '~/composables/useStore'

const route = useRoute()
const router = useRouter()
const store = useStore()

const username = ref('')
const password = ref('')

const login = async () => {
  const resp = await reqAuthenticate({
    username: username.value,
    password: password.value,
  })

  if (isError(resp)) {
    store.setToast(resp.error)
    return
  }

  tokenSet(resp.result.token)
  router.replace((route.query.redirect as string) || '/')
}
</script>
