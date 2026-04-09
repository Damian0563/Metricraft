<template>
	<NuxtLayout>
		<Navbar />
		<Spinner :loading="loading" />
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Sign :oldUser="oldUserStatus" @signup="handleSignup" @load="handleLoad" @popup="createPop"
			@toggle="oldUserStatus = !oldUserStatus" />
	</NuxtLayout>
</template>


<script setup lang="ts">
import { welcome } from "~/calls/welcome"
const oldUserStatus = ref(false)
const errorMessage = ref("")
const loading = ref(true)

const handleStatus = async () => {
	try {
		const result: boolean | null = await welcome()
		if (result == null) {
			errorMessage.value = "Something went wrong"
			return
		}
		oldUserStatus.value = result
		if (oldUserStatus.value) {
			await navigateTo("/dashboard")
		}
	} catch (error) {
		errorMessage.value = "Something went wrong, Check your internet connection and try again."
	}
}

const handleLoad = () => {
	loading.value = !loading.value
}

const handleSignup = async (uuid: string) => {
	loading.value = true
	document.cookie = `session-token=${uuid}; path=/;expires=${new Date(Date.now() + 3600000).toUTCString()};SameSite=None;Secure`
	await navigateTo("/dashboard")
	loading.value = false
}

const createPop = (msg: string) => {
	errorMessage.value = msg
}

onMounted(async () => {
	await handleStatus()
	loading.value = false
})
</script>
