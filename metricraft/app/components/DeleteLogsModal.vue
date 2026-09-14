<template>
	<ClientOnly>
		<Teleport to="body">
			<AnimatePresence>
				<motion.div v-if="show" key="delete-logs-overlay"
					class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm md:p-8"
					:initial="{ opacity: 0 }" :animate="{ opacity: 1 }" :exit="{ opacity: 0 }" :transition="{ duration: 0.2 }"
					@click.self="close">
					<motion.div class="flex w-full max-w-md flex-col overflow-hidden rounded-xl bg-white shadow-xl ring-1 ring-slate-100"
						role="dialog" aria-modal="true" :initial="{ opacity: 0, scale: 0.96, y: 16 }"
						:animate="{ opacity: 1, scale: 1, y: 0 }" :exit="{ opacity: 0, scale: 0.96, y: 16 }"
						:transition="{ type: 'spring', duration: 0.35, bounce: 0.2 }" @click.stop>
						<div class="flex items-center justify-between px-6 py-4 border-b border-gray-100">
							<h2 class="text-lg font-semibold text-gray-800">Free up log storage</h2>
							<button type="button" @click="close"
								class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-dark-gray transition-colors hover:bg-slate-100 hover:text-[#00F376]"
								aria-label="Close">
								<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
									stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5">
									<line x1="18" y1="6" x2="6" y2="18"></line>
									<line x1="6" y1="6" x2="18" y2="18"></line>
								</svg>
							</button>
						</div>
						<div class="space-y-5 px-6 py-5">
							<p class="text-sm text-gray-600">
								The oldest logs are deleted first until roughly the selected amount of storage is freed.
								This cannot be undone.
							</p>
							<div class="flex items-center justify-between text-sm">
								<span class="text-gray-500">Current storage</span>
								<span class="font-semibold text-gray-900">{{ props.capacity.toFixed(2) }} MB</span>
							</div>
							<div class="space-y-3">
								<label for="delete-mb" class="text-sm font-medium text-gray-700">Amount to delete</label>
								<input type="range" min="0" :max="props.capacity" step="0.01" v-model.number="toDeleteMbs"
									class="w-full accent-[#00F376]" />
								<div class="flex items-center gap-2">
									<input id="delete-mb" type="number" min="0" :max="props.capacity" step="0.01"
										v-model.number="toDeleteMbs"
										class="w-32 px-3 py-2 rounded-lg border border-gray-200 text-sm text-gray-800 focus:outline-none focus:border-[#00F376]" />
									<span class="text-sm text-gray-500">MB</span>
									<button type="button" @click="toDeleteMbs = props.capacity"
										class="ml-auto text-xs font-medium text-gray-500 hover:text-[#00F376] transition-colors">
										Select all
									</button>
								</div>
								<p class="text-xs text-gray-500">
									About {{ remaining.toFixed(2) }} MB will remain.
								</p>
							</div>
						</div>
						<div class="flex justify-end gap-3 px-6 py-4 border-t border-gray-100">
							<button type="button" @click="close"
								class="px-4 py-2 rounded-lg text-sm font-medium text-gray-700 hover:bg-slate-100 transition-colors">
								Cancel
							</button>
							<button type="button" @click="confirm" :disabled="!isValid"
								class="px-4 py-2 rounded-lg text-sm font-semibold text-white bg-red-500 hover:bg-red-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
								Delete {{ isValid ? `${toDeleteMbs.toFixed(2)} MB` : 'logs' }}
							</button>
						</div>
					</motion.div>
				</motion.div>
			</AnimatePresence>
		</Teleport>
	</ClientOnly>
</template>

<script setup lang="ts">
import { motion, AnimatePresence } from 'motion-v';

const props = defineProps<{
	show: boolean;
	capacity: number;
}>();
const emit = defineEmits<{
	close: [];
	confirm: [mb: number];
}>();
const toDeleteMbs = ref(0)
const isValid = computed(() => typeof toDeleteMbs.value === 'number' && toDeleteMbs.value > 0 && toDeleteMbs.value <= props.capacity)
const remaining = computed(() => Math.max(props.capacity - (Number(toDeleteMbs.value) || 0), 0))

watch(() => props.show, (visible) => {
	if (visible) toDeleteMbs.value = 0
})

const close = () => emit('close')
const confirm = () => {
	if (!isValid.value) return
	emit('confirm', toDeleteMbs.value)
}
</script>
