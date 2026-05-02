<template>
	<div class="max-w-6xl mx-auto px-16 py-8">
		<div class="flex items-center justify-between mb-12">
			<h1 class="text-2xl font-bold text-black">Settings</h1>
			<button @click="emit('back')"
				class="hover:cursor-pointer p-2 rounded-full transition-all duration-200 hover:bg-gray-200">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-700" fill="none" viewBox="0 0 24 24"
					stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
				</svg>
			</button>
		</div>
		<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
			<div class="space-y-8">
				<div class="pb-6 border-b border-gray-200">
					<h2 class="text-xl font-semibold text-gray-800">Application</h2>
					<p class="text-sm text-gray-500 mt-2">{{ appName }}</p>
				</div>
				<div class="pb-6 border-b border-gray-200">
					<h2 class="text-xl font-semibold text-gray-800 mb-4">Preferences</h2>
					<label class="flex items-center justify-between cursor-pointer">
						<span class="text-base font-medium text-gray-700">Real-time updates</span>
						<div class="relative">
							<input type="checkbox" class="sr-only peer" :checked="realtimeEnabled"
								@change="emit('realtimeToggle', ($event.target as HTMLInputElement).checked)" />
							<div
								class="w-11 h-6 bg-gray-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#00F376]">
							</div>
						</div>
					</label>
				</div>
				<div class="pb-6 border-b border-gray-200">
					<h2 class="text-xl font-semibold text-gray-800 mb-4">Customization</h2>
					<button @click="emit('customizeView', !customizeDashboard)"
						class="w-full text-left px-4 py-3 rounded-xl border border-gray-200 hover:border-[#00F376] hover:shadow-md transition-all duration-300 text-gray-700 font-medium">
						Customize Dashboard View
					</button>
				</div>
				<div>
					<h2 class="text-xl font-semibold text-gray-800 mb-4">Team</h2>
					<span @click="copyInvite"
						class="text-base font-medium text-gray-700 hover:cursor-pointer hover:text-[#00F376] transition-colors duration-200">
						Invite team members
					</span>
				</div>
			</div>
		</div>
		<div class="mt-8 flex justify-end">
			<button @click="emit('load')"
				class="bg-[#00F376] text-black rounded-full px-6 py-3 text-sm font-medium transition-all duration-300 shadow-md hover:shadow-lg hover:cursor-pointer">
				Save Changes
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
const props = defineProps<{
	appName: string;
	realtimeEnabled: boolean;
}>();
const emit = defineEmits<{
	realtimeToggle: [value: boolean];
	customizeView: [value: boolean];
	load: [value: void];
	back: [value: void];
}>();
const customizeDashboard = ref(false)
const copyInvite = () => {
	emit('load')
	const cookie: string | undefined = document.cookie.split(";").find(c => c.trim().startsWith("session-token"))?.split("=")[1]
	if (!cookie) return
	navigator.clipboard.writeText(`http://localhost:8000/invite/inv:${cookie}`)
	emit('load')
}
</script>
