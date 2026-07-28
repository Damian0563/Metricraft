<template>
	<ClientOnly>
		<Teleport to="body">
			<div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm px-4">
				<div class="relative bg-white p-8 shadow-xl rounded-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
					<h2 class="text-xl font-semibold text-gray-800 mb-2">Edit worker</h2>
					<p class="text-sm text-gray-500 mb-6">
						Update the upstream URL, poll interval, and headers for this worker.
					</p>
					<WorkerForm :worker="worker" :saving="saving" show-exit id-prefix="edit-worker" @close="emit('close')"
						@save="emit('save', $event)" />
				</div>
			</div>
		</Teleport>
	</ClientOnly>
</template>

<script setup lang="ts">
import type { Worker } from '@/composables/types/additional'

defineProps<{
	open: boolean
	worker: Worker | null
	saving?: boolean
}>()

const emit = defineEmits<{
	close: []
	save: [worker: Worker]
}>()
</script>
