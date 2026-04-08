<template>
	<NuxtLayout>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Spinner :loading="loading" />
	</NuxtLayout>
</template>

<script setup lang="ts">
import { getCookie, updateCookie } from "@/composables/helpers";
import type { dashboardInitPayload } from '@/composables/types';
import { getDashboard } from "~/calls/dashboard";
let cookie = getCookie("session-token");
const loading = ref(true)
const errorMessage = ref("")
const appName = ref("")

const init = async () => {
	loading.value = true
	try {
		const data: dashboardInitPayload = await getDashboard(cookie)
		if (data.error) {
			errorMessage.value = data.error
		} else {
			appName.value = data.appName
			cookie = data.signedSecret
			updateCookie(cookie)
		}
	} catch (e) {
		errorMessage.value = e as string
	} finally {
		loading.value = false
	}
}

onMounted(() => {
	init()
})
</script>
