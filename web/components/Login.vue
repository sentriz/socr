<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-200">
    <div class="m-8 w-full max-w-xs space-y-4">
      <div class="space-y-6 rounded bg-white p-8 shadow-md">
        <logo class="mx-auto w-9/12" />
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
      <p class="text-center text-xs text-gray-500"><b>s</b>creenshot <b>ocr</b> server &mdash; Senan Kelly 2020</p>
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
