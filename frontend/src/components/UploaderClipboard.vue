<template>
  <loading-modal :loading="loading" text="uploading from clipboard" />
</template>

<script setup lang="ts">
import LoadingModal from './LoadingModal.vue'

import { isError, reqUpload } from '../api'
import { useRouter } from 'vue-router'
import useLoading from '../composables/useLoading'

const router = useRouter()
const { loading, load } = useLoading(reqUpload)

// TODO: not attach handler to the whole document, just the main page
document.onpaste = async (event: ClipboardEvent) => {
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
