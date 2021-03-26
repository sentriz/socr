<template>
  <TransitionFade>
    <div
      @paste="paste"
      v-if="loading"
      class="z-10 fixed inset-0 bg-gray-700 bg-opacity-75 transition-opacity flex items-center justify-center"
    >
      <LoadingSpinner text="uploading" />
    </div>
  </TransitionFade>
</template>

<script setup lang="ts">
import LoadingSpinner from './LoadingSpinner.vue'
import TransitionFade from './TransitionFade.vue'

import { isError, reqUpload } from '../api'
import { useRouter } from 'vue-router'
import useLoading from '../composables/useLoading'

const router = useRouter()

const { loading, load } = useLoading(reqUpload)

const paste = async (event: ClipboardEvent) => {
  const items = event.clipboardData?.items
  if (!items) return
  const blob = items[0]?.getAsFile()
  if (!blob) return

  const formData = new FormData()
  formData.append('i', blob)

  const resp = await load(formData)
  if (isError(resp)) return

  router.push({ name: 'public', params: { hash: resp.result.id } })
}
</script>
