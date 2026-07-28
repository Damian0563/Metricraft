<template>
	<ClientOnly>
		<Teleport to="body">
			<div v-if="props.message"
				class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm px-4">
				<div class="relative bg-white p-8 shadow-xl rounded-xl w-full max-w-md">
					<button type="button" @click="emit('close')" aria-label="Close"
						class="absolute top-3 right-3 w-8 h-8 flex items-center justify-center rounded-full text-gray-400 hover:text-gray-700 hover:bg-gray-100 transition">
						<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
							stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5">
							<line x1="18" y1="6" x2="6" y2="18" />
							<line x1="6" y1="6" x2="18" y2="18" />
						</svg>
					</button>
					<div class="flex justify-center mb-4">
						<div class="w-14 h-14 rounded-full bg-[#00F376]/15 flex items-center justify-center">
							<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#00F376" stroke-width="2"
								stroke-linecap="round" stroke-linejoin="round" class="w-7 h-7">
								<path d="M12 9v4m0 4h.01M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" />
							</svg>
						</div>
					</div>
					<h2 class="text-2xl font-bold text-center text-gray-800 mb-2">{{ props.title }}</h2>
					<p class="text-sm text-gray-500 text-center mb-6">
						{{ props.message }}
					</p>
					<button type="button" @click="emit('close')"
						class="block mx-auto px-6 py-2 text-black font-semibold bg-[#00F376] hover:text-white rounded-lg shadow-lg hover:bg-black transition delay-100 ease-in-out">
						Got it
					</button>
				</div>
			</div>
		</Teleport>
	</ClientOnly>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
	message: string;
	title?: string;
}>(), {
	title: 'Heads up',
});
const emit = defineEmits<{
	close: [];
}>();
</script>
