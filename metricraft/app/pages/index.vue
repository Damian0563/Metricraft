<template>
	<div class="bg-black h-screen">
		<Navbar />
		<Signup :newUser="newUserStatus" @signup="handleSignup" @load="handleLoad" @popup="createPop" />
	</div>
</template>


<script setup lang="ts">
import { welcome } from "@/helpers/welcome"
const newUserStatus = ref(false)
const errorMessage = ref("")
const loading = ref(true)

const handleStatus = async () => {
	try {
		const result: boolean | null = await welcome()
		if (result == null) {
			errorMessage.value = "Something went wrong"
			return
		}
		//newUserStatus.value = result
		newUserStatus.value = true
	} catch (error) {
		errorMessage.value = error as string
	}
}

const handleLoad = () => {
	loading.value = !loading.value
}

const handleSignup = () => {
	loading.value = true
}

const createPop = (msg: string) => {
	console.log("create pop", msg)
}

onMounted(async () => {
	await handleStatus()
	loading.value = false
})
</script>
