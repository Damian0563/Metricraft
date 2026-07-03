<template>
	<div class="text-white justify-start text-2xl font-bold py-4 px-16 mb-8" style="background-color: #00F376;">
		<nav>
			<div class="flex justify-between">
				<div class="flex items-center">
					<img src="/logo.svg" class="h-8 mr-4 cursor-pointer" alt="logo, go back to main view" @click="getBack" />
					<h1 class="text-black text-2xl font-bold">{{ appName }}</h1>
				</div>
				<div class="flex items-center gap-8">
					<button
						class="nav-menu-btn bg-white text-black rounded-full px-4 py-2 text-sm hover:cursor-pointer transition-all duration-300"
						@click="options = !options">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
							stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
						</svg>
					</button>
					<div class="relative">
						<div v-if="options"
							class="dropdown-panel absolute top-0 right-0 z-10 mt-6 mr-2 bg-white rounded-xl p-3 min-w-48 border border-gray-100 flex flex-col gap-1">
							<button type="button" class="dropdown-item" @click="emit('team'); options = false">
								Team
							</button>
							<button type="button" class="dropdown-item" @click="emit('workers'); options = false">
								Metricraft workers
							</button>
							<button type="button" class="dropdown-item" @click="emit('settings'); options = false">
								Settings
							</button>
							<button
								class="sign-out-btn text-sm font-medium text-white bg-red-500 rounded-full px-3 py-2 w-full mt-2 transition-all duration-300 hover:bg-red-600"
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
	workers: [value: void];
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

<style scoped>
.dropdown-panel {
	box-shadow:
		0 4px 6px rgba(0, 0, 0, 0.04),
		0 10px 28px rgba(0, 0, 0, 0.1);
}

.dropdown-item {
	display: block;
	width: 100%;
	padding: 0.5rem 0.75rem;
	text-align: center;
	font-size: 0.875rem;
	font-weight: 500;
	color: rgb(55 65 81);
	border-radius: 0.5rem;
	cursor: pointer;
	background: transparent;
	border: none;
	transition: background-color 0.2s ease, color 0.2s ease;
}

.dropdown-item:hover {
	background-color: rgb(249 250 251);
	color: rgb(17 24 39);
}

.dropdown-item:active {
	background-color: rgb(243 244 246);
}

.nav-menu-btn {
	box-shadow:
		0 1px 2px rgba(0, 0, 0, 0.08),
		0 4px 14px rgba(0, 0, 0, 0.12);
}

.nav-menu-btn:hover {
	box-shadow:
		0 2px 4px rgba(0, 0, 0, 0.1),
		0 8px 24px rgba(0, 0, 0, 0.16);
}

.nav-menu-btn:active {
	box-shadow:
		0 1px 2px rgba(0, 0, 0, 0.1),
		0 2px 6px rgba(0, 0, 0, 0.12);
}

.sign-out-btn {
	box-shadow:
		0 2px 6px rgba(239, 68, 68, 0.35),
		0 1px 3px rgba(185, 28, 28, 0.2);
}

.sign-out-btn:hover {
	box-shadow:
		0 4px 14px rgba(239, 68, 68, 0.45),
		0 2px 6px rgba(185, 28, 28, 0.25);
}

.sign-out-btn:active {
	box-shadow:
		0 1px 3px rgba(239, 68, 68, 0.3),
		0 1px 2px rgba(185, 28, 28, 0.2);
}
</style>
