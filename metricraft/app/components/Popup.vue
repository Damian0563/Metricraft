<template>
	<div class="fixed bottom-4 right-4 bg-white px-4 py-5 shadow-xl rounded-xl w-64 z-50" v-if="props.message">
		<button type="button" @click="emit('close')"
			class="absolute top-2 right-2 w-6 h-6 flex items-center justify-center text-[#00F376] hover:text-[#00d466] transition-colors">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
				stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5">
				<line x1="18" y1="6" x2="6" y2="18"></line>
				<line x1="6" y1="6" x2="18" y2="18"></line>
			</svg>
		</button>
		<div class="flex items-start gap-3 pr-6">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
				stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 flex-shrink-0 mt-0.5">
				<path
					d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
			</svg>
			<p class="text-sm text-gray-700">{{ props.message }}</p>
		</div>
		<div class="mt-3 h-2 w-full bg-gray-200 rounded-full overflow-hidden">
			<div class="h-full bg-[#00F376] rounded-full transition-all duration-50 ease-linear"
				:style="{ width: `${progress}%` }"></div>
		</div>
	</div>
</template>

<script setup lang="ts">
const emit = defineEmits(["close"])
const progress = ref(0)
let intervalId: ReturnType<typeof setInterval> | null = null
const props = defineProps({
	message: {
		type: String,
		required: true,
	},
})
watch(() => props.message, (newMessage) => {
	if (!newMessage) return
	if (intervalId) clearInterval(intervalId)
	progress.value = 0
	const duration = 5000
	const startTime = Date.now()
	intervalId = setInterval(() => {
		const elapsed = Date.now() - startTime
		progress.value = Math.min((elapsed / duration) * 100, 100)
		if (progress.value >= 100) {
			clearInterval(intervalId!)
			emit("close")
		}
	}, 50)
}, { immediate: true })
</script>
