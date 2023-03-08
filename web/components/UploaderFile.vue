<template>
  <input class="hidden" ref="elm" type="file" accept="image/*, video/*" @change="select" />
  <div class="padded pointer-events-auto absolute bottom-0 right-0 rounded-tl bg-gray-200 opacity-75 hover:opacity-100">
    <div class="flex items-center gap-2 text-gray-600 hover:cursor-pointer hover:text-gray-900" @click="click">
      <upload-icon class="h-5" />
      <span>upload or paste file</span>
    </div>
  </div>
  <loading-modal :loading="loading" text="uploading from file" />
</template>

<script setup lang="ts">
import LoadingModal from './LoadingModal.vue'

import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { isError, reqUpload } from '~/request'
import useLoading from '~/composables/useLoading'
import { UploadIcon } from '@heroicons/vue/outline'
import { routes } from '~/router'

const router = useRouter()
const { loading, load } = useLoading(reqUpload)

const elm = ref<HTMLInputElement>()

const click = () => elm.value?.click()
const select = async () => {
  const file = elm.value?.files?.[0]
  if (!file) return

  const formData = new FormData()
  formData.append('i', file)

  const resp = await load(formData)
  if (isError(resp)) return

  router.push({ name: routes.PUBLIC, params: { hash: resp.result.id } })
}
</script>
