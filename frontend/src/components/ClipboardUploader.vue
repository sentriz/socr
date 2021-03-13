<template></template>

<script setup lang="ts">
import { isError, reqUpload } from '../api'

document.onpaste = async (event) => {
  const items = event.clipboardData?.items
  if (!items) return

  const blob = items[0]?.getAsFile()
  if (!blob) return

  const formData = new FormData()
  formData.append('i', blob)

  const r = await reqUpload(formData)
  if (isError(r)) return

  console.log(r)
}
</script>
