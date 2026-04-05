<template>
	<div class="bg-black" :class="overflow ? 'h-auto' : 'h-screen'">
		<Navbar />
		<Spinner :loading="loading" />
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Sign :oldUser="oldUserStatus" @signup="handleSignup" @load="handleLoad" @popup="createPop"
			@toggle="oldUserStatus = !oldUserStatus" />
	</div>
</template>


<script setup lang="ts">
import { welcome } from "~/calls/welcome"
const oldUserStatus = ref(false)
const errorMessage = ref("")
const loading = ref(true)
const overflow = ref(false)
const checkOverflow = () => {
	overflow.value = document.documentElement.scrollHeight > window.innerHeight
}
watch(oldUserStatus, () => {
	nextTick(checkOverflow)
})
onMounted(() => {
	checkOverflow()
	window.addEventListener('resize', checkOverflow)
})
onUnmounted(() => {
	window.removeEventListener('resize', checkOverflow)
})
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
		errorMessage.value = error as string
	}
}

const handleLoad = () => {
	loading.value = !loading.value
}

const handleSignup = async (uuid: string) => {
	loading.value = true
	document.cookie = `session-token=${uuid}; path=/`
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
