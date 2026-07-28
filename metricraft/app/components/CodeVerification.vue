<template>
	<ClientOnly>
		<Teleport to="body">
			<div v-if="props.status"
				class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm px-4">
				<div class="relative bg-white p-8 shadow-xl rounded-xl w-full max-w-md">
					<button type="button" @click="closeModal" aria-label="Close"
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
								<path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" />
								<polyline points="22,6 12,13 2,6" />
							</svg>
						</div>
					</div>
					<h2 class="text-2xl font-bold text-center text-gray-800 mb-2">Verify your email</h2>

					<p class="text-sm text-gray-500 text-center mb-6">
						Enter the 6-digit code we just sent to your inbox.
					</p>
					<div class="flex justify-center gap-2 mb-6" @paste.prevent="handlePaste">
						<input v-for="(_, idx) in 6" :key="idx" :ref="(el) => setInputRef(el, idx)" type="text" inputmode="numeric"
							autocomplete="one-time-code" maxlength="1" v-model="digits[idx]" @input="handleInput($event, idx)"
							@keydown="handleKeydown($event, idx)" @focus="handleFocus(idx)" :disabled="locked"
							class="w-11 h-14 text-center text-xl font-semibold text-gray-800 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed" />
					</div>
					<button type="button" @click="resetCode" :disabled="locked"
						class="block mx-auto text-sm font-medium text-gray-500 hover:text-[#00F376] transition disabled:opacity-50 disabled:hover:text-gray-500">
						Clear code
					</button>
				</div>
			</div>
		</Teleport>
	</ClientOnly>
</template>

<script setup lang="ts">
const props = defineProps<{
	status: boolean;
}>();
const emit = defineEmits<{
	complete: [value: string];
}>();

const digits = ref<string[]>(Array(6).fill(''));
const inputRefs = ref<HTMLInputElement[]>([]);
const locked = ref(false);
const setInputRef = (el: unknown, idx: number) => {
	if (el) inputRefs.value[idx] = el as HTMLInputElement;
};
const code = computed(() => digits.value.join(''));
watch(code, (val) => {
	if (val.length === 6 && /^\d{6}$/.test(val)) {
		locked.value = true;
		emit('complete', val);
	}
});
const handleInput = (e: Event, idx: number) => {
	const target = e.target as HTMLInputElement;
	const value = target.value.replace(/\D/g, '');
	digits.value[idx] = value.slice(-1);
	if (digits.value[idx] && idx < 5) {
		nextTick(() => inputRefs.value[idx + 1]?.focus());
	}
};

const handleKeydown = (e: KeyboardEvent, idx: number) => {
	if (e.key === 'Backspace') {
		if (!digits.value[idx] && idx > 0) {
			e.preventDefault();
			inputRefs.value[idx - 1]?.focus();
			digits.value[idx - 1] = '';
		} else {
			digits.value[idx] = '';
		}
	} else if (e.key === 'ArrowLeft' && idx > 0) {
		e.preventDefault();
		inputRefs.value[idx - 1]?.focus();
	} else if (e.key === 'ArrowRight' && idx < 5) {
		e.preventDefault();
		inputRefs.value[idx + 1]?.focus();
	}
};

const handleFocus = (idx: number) => {
	inputRefs.value[idx]?.select();
};

const handlePaste = (e: ClipboardEvent) => {
	const pasted = e.clipboardData?.getData('text').replace(/\D/g, '').slice(0, 6) ?? '';
	if (!pasted) return;
	for (let i = 0; i < 6; i++) {
		digits.value[i] = pasted[i] ?? '';
	}
	const nextIdx = Math.min(pasted.length, 5);
	nextTick(() => inputRefs.value[nextIdx]?.focus());
};

const resetCode = () => {
	digits.value = Array(6).fill('');
	nextTick(() => inputRefs.value[0]?.focus());
};

const closeModal = () => {
	emit('complete', 'CLOSED');
};

watch(() => props.status, async (val) => {
	if (val) {
		digits.value = Array(6).fill('');
		locked.value = false;
		await nextTick();
		inputRefs.value[0]?.focus();
	}
});
</script>
