<template>
	<div class="bg-black h-screen">
		<Navbar />
		<Signup />
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
		newUserStatus.value = result
	} catch (error) {
		errorMessage.value = error as string
	}
}


onMounted(async () => {
	await handleStatus()
	loading.value = false
})
</script>
