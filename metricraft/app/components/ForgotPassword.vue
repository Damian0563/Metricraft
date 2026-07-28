<template>
	<ClientOnly>
		<Teleport to="body">
			<div v-if="props.status"
				class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm px-4">
				<div class="relative bg-white p-8 shadow-xl rounded-xl w-full max-w-md">
					<button type="button" @click="closeModal" aria-label="Close"
						class="absolute top-3 right-3 w-8 h-8 flex items-center justify-center rounded-full text-gray-400 hover:text-gray-700 hover:bg-gray-100 transition">
						<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
							stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5">
							<line x1="18" y1="6" x2="6" y2="18" />
							<line x1="6" y1="6" x2="18" y2="18" />
						</svg>
					</button>
					<div class="flex justify-center mb-4">
						<div class="w-14 h-14 rounded-full bg-[#00F376]/15 flex items-center justify-center">
							<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#00F376" stroke-width="2"
								stroke-linecap="round" stroke-linejoin="round" class="w-7 h-7">
								<rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
								<path d="M7 11V7a5 5 0 0 1 10 0v4" />
							</svg>
						</div>
					</div>
					<h2 class="text-2xl font-bold text-center text-gray-800 mb-2">Forgot your password?</h2>
					<p class="text-sm text-gray-500 text-center mb-6">
						Enter your email and we'll send you a recovery link.
					</p>
					<form @submit.prevent="handleSubmit" class="flex flex-col gap-4">
						<div>
							<label for="recoveryMail" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
							<input id="recoveryMail" ref="inputRef" type="email" v-model="mail" placeholder="you@example.com"
								:disabled="loading"
								class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed" />
						</div>
						<button type="submit" :disabled="loading"
							class="w-full cursor-pointer py-3 mt-2 text-black font-semibold bg-[#00F376] hover:text-white rounded-lg shadow-lg hover:bg-black transition delay-100 ease-in-out disabled:opacity-60 disabled:cursor-not-allowed disabled:hover:bg-[#00F376] disabled:hover:text-black">
							{{ loading ? 'Sending...' : 'Send' }}
						</button>
					</form>
				</div>
			</div>
		</Teleport>
	</ClientOnly>
</template>

<script setup lang="ts">
const props = defineProps<{
	status: boolean;
	loading?: boolean;
}>();
const emit = defineEmits<{
	submit: [value: string];
	close: [];
}>();
const mail = ref('');
const inputRef = ref<HTMLInputElement | null>(null);
const closeModal = () => {
	emit('close');
};
const handleSubmit = () => {
	emit('submit', mail.value);
};
watch(() => props.status, async (val) => {
	if (val) {
		mail.value = '';
		await nextTick();
		inputRef.value?.focus();
	}
});
</script>
