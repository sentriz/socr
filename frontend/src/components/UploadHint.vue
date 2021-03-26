<template>
  <input class="hidden" ref="elm" type="file" accept="image/*" @change="select" />
  <div class="fixed bottom-0 opacity-75 hover:opacity-100 right-0 padded bg-gray-200 rounded-tl">
    <div class="flex items-center gap-2 text-gray-600 hover:text-gray-900 hover:cursor-pointer" @click="click">
      <i class="fas fa-upload text-sm"></i>
      <span>select file to upload or just paste</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { isError, reqUpload } from '../api'

const router = useRouter()

const elm = ref<HTMLInputElement>()

const click = () => elm.value?.click()
const select = async () => {
  const file = elm.value?.files?.[0]
  if (!file) return

  const formData = new FormData()
  formData.append('i', file)

  const resp = await reqUpload(formData)
  if (isError(resp)) return

  router.push({ name: 'public', params: { hash: resp.result.id } })
}
</script>
