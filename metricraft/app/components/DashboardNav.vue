<template>
	<div v-if="props.appName" class="text-white justify-start text-2xl font-bold py-8 px-16"
		style="background-color: #00F376;">
		<nav>
			<div class="flex justify-between">
				<div class="flex items-center">
					<img src="/logo.svg" class="h-8 mr-4" alt="logo" />
					<h1 class="text-black text-2xl font-bold">{{ props.appName }}</h1>
				</div>
				<div class="flex items-center">
					<button class="bg-white text-black rounded-full px-4 py-2 text-sm" @click="options = !options">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
							stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
						</svg>
					</button>
					<div class="relative" v-if="options">
						<div
							class="absolute top-0 right-0 z-10 mt-6 mr-2 bg-white rounded-xl shadow-xl p-4 min-w-48 border border-gray-100">
							<label class="flex items-center gap-3 cursor-pointer">
								<span class="text-sm font-medium text-gray-700">
									Real-time updates
								</span>
								<input type="checkbox"
									class="w-4 h-4 rounded border-gray-300 text-[#00F376] focus:ring-[#00F376] cursor-pointer"
									@toggle="emit('realTimeToggle', $event.target.value)" />
							</label>
						</div>
					</div>
				</div>
			</div>
		</nav>
	</div>
	<div class="flex flex-col p-16 mx-8">
	</div>

</template>

<script setup lang="ts">
const props = defineProps<{
	appName: string;
}>();
const emit = defineEmits<{
	realTimeToggle: boolean;
}>();
const options = ref(false)

</script>
