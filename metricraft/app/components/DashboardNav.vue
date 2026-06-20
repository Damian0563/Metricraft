<template>
	<div class="text-white justify-start text-2xl font-bold py-4 px-16 mb-8" style="background-color: #00F376;">
		<nav>
			<div class="flex justify-between">
				<div class="flex items-center">
					<img src="/logo.svg" class="h-8 mr-4 cursor-pointer" alt="logo, go back to main view" @click="getBack" />
					<h1 class="text-black text-2xl font-bold">{{ appName }}</h1>
				</div>
				<div class="flex items-center gap-8">
					<button class="bg-white text-black rounded-full px-4 py-2 text-sm hover:cursor-pointer"
						@click="options = !options">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
							stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
						</svg>
					</button>
					<div class="relative">
						<div v-if="options"
							class="absolute top-0 right-0 z-10 mt-6 mr-2 bg-white rounded-xl shadow-xl p-4 min-w-48 border border-gray-100 justfy-center">
							<span @click="emit('team'); options = false"
								class="text-sm font-medium text-gray-700 mt-4 hover:cursor-pointer text-center block">
								Team
							</span>
							<span @click="emit('settings'); options = false"
								class="text-sm font-medium text-gray-700 mt-4 hover:cursor-pointer text-center block">
								Settings
							</span>
							<button
								class="text-sm font-medium text-white bg-red-500 rounded-full p-1 w-full transition-all duration-300 shadow-md hover:bg-red-600 hover:shadow-lg"
								@click="signOut">Sign out
							</button>
						</div>
					</div>
				</div>
			</div>
		</nav>
	</div>
</template>

<script setup lang="ts">
import { invalidateCookie } from "~/composables/helpers"
const appName = useState<string>('appName');
const emit = defineEmits<{
	settings: [value: void];
	team: [value: void];
}>();
const options = ref(false)
const signOut = () => {
	invalidateCookie()
	navigateTo('/')
}
const router = useRouter()
const getBack = () => {
	if (router.currentRoute.value.name === 'dashboard') {
		router.replace({ query: {} })
	} else {
		navigateTo('/dashboard')
	}
}
</script>
