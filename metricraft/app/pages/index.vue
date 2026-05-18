<template>
	<div>
		<Navbar />
		<ClientOnly>
			<Spinner :loading="loading || localLoading" />
		</ClientOnly>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Sign :oldUser="mode" @signup="handleSignup" @load="handleLoad" @popup="createPop" @toggle="mode = !mode" />
	</div>
</template>


<script setup lang="ts">
definePageMeta({
	layout: 'entry',
})
import { welcome } from "~/calls/welcome"
const errorMessage = ref("")
const mode = ref(true)
const localLoading = ref(false)
const { data: oldUserStatus, pending: loading, error } = await useAsyncData('welcome', () => welcome())
if (oldUserStatus.value) {
	navigateTo("/dashboard")
}
watch(oldUserStatus, (newVal) => {
	if (newVal) {
		navigateTo("/dashboard")
	}
})
watch(loading, (newVal) => {
	localLoading.value = newVal
})
watch(error, (newVal) => {
	if (newVal) {
		errorMessage.value = "Something went wrong, Check your internet connection and try again."
	}
})
const handleLoad = () => {
	localLoading.value = !localLoading.value
}
const handleSignup = async (uuid: string) => {
	localLoading.value = true
	document.cookie = `session-token=${uuid}; path=/;expires=${new Date(Date.now() + 3600000 * 24).toUTCString()};SameSite=None;Secure`
	await navigateTo("/dashboard")
	localLoading.value = false
}

const createPop = (msg: string) => {
	errorMessage.value = msg
}
</script>
