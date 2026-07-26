<template>
	<div class="space-y-6">
		<div>
			<label :for="`${idPrefix}-url`" class="block text-sm font-medium text-gray-700 mb-2">Upstream URL</label>
			<input :id="`${idPrefix}-url`" v-model="editUrl" type="url" placeholder="https://api.example.com/health"
				:disabled="showExit"
				class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors disabled:bg-gray-100 disabled:cursor-not-allowed" />
			<p v-if="showHelperText" class="mt-2 text-xs text-gray-500">
				The full address of the service the worker proxies traffic to and records metrics for.
			</p>
		</div>

		<div class="flex flex-col gap-8 xl:flex-row xl:justify-between">
			<div class="shrink-0">
				<label :for="`${idPrefix}-poll`" class="block text-sm font-medium text-gray-700 mb-2">Poll interval
					(minutes)</label>
				<input :id="`${idPrefix}-poll`" v-model.number="editPollInterval" type="number" min="5" step="1" max="60"
					placeholder="10" :disabled="saving"
					class="w-full max-w-xs px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors disabled:bg-gray-100 disabled:cursor-not-allowed" />
				<p v-if="showHelperText" class="mt-2 text-xs text-gray-500 max-w-xs">
					How often the worker sends a health check request to keep metrics fresh between real user traffic.
				</p>
			</div>
			<div class="flex-1 min-w-0">
				<div class="flex items-center justify-between mb-2">
					<label class="block text-sm font-medium text-gray-700">Headers</label>
					<button type="button" @click="addHeader" :disabled="saving"
						class="text-sm font-medium text-[#00F376] hover:text-[#00D96A] transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed">
						+ Add header
					</button>
				</div>

				<div v-if="headerRows.length === 0"
					class="rounded-lg border border-dashed border-gray-200 px-4 py-6 text-center">
					<p class="text-sm text-gray-500">No custom headers configured.</p>
					<button type="button" @click="addHeader" :disabled="saving"
						class="mt-2 text-sm font-medium text-[#00F376] hover:text-[#00D96A] transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed">
						Add your first header
					</button>
				</div>

				<div v-else class="space-y-2 max-h-48 overflow-y-auto pr-1">
					<div v-for="(header, index) in headerRows" :key="index" class="flex items-center gap-2">
						<input v-model="header.key" type="text" placeholder="Authorization" :disabled="saving"
							class="w-2/5 min-w-0 px-3 py-2 rounded-lg border border-gray-200 text-gray-800 text-sm focus:outline-none focus:border-[#00F376] transition-colors disabled:bg-gray-100 disabled:cursor-not-allowed" />
						<input v-model="header.value" type="text" placeholder="Bearer token…" :disabled="saving"
							class="flex-1 min-w-0 px-3 py-2 rounded-lg border border-gray-200 text-gray-800 text-sm focus:outline-none focus:border-[#00F376] transition-colors disabled:bg-gray-100 disabled:cursor-not-allowed" />
						<button type="button" @click="removeHeader(index)" :disabled="saving"
							class="shrink-0 p-2 text-gray-400 hover:text-red-500 transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
							aria-label="Remove header">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd"
									d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
									clip-rule="evenodd" />
							</svg>
						</button>
					</div>
				</div>

				<p v-if="showHelperText" class="mt-2 text-xs text-gray-500">
					Optional HTTP headers included on each poll request, for example auth tokens or custom
					<code class="text-gray-600">Accept</code> values.
				</p>
			</div>
		</div>

		<div class="pt-6 border-t border-gray-100 flex justify-end gap-3">
			<button v-if="showExit" type="button" @click="emit('close')" :disabled="saving"
				class="px-6 py-3 bg-gray-100 text-gray-700 font-semibold rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer">
				Exit
			</button>
			<button type="button" @click="handleSave" :disabled="!canSave || saving"
				class="px-8 py-3 bg-[#00F376] text-gray-900 font-bold rounded-lg hover:bg-[#00D96A] transition-all shadow-lg disabled:opacity-50 disabled:cursor-not-allowed uppercase tracking-wider text-sm cursor-pointer">
				{{ saveLabel }}
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import type { Worker } from '@/composables/types/additional'
import { useWorkerForm } from '@/composables/workers'

const props = withDefaults(defineProps<{
	worker?: Worker | null
	saving?: boolean
	saveLabel?: string
	showHelperText?: boolean
	showExit?: boolean
	idPrefix?: string
}>(), {
	worker: null,
	saving: false,
	saveLabel: 'Save',
	showHelperText: false,
	showExit: false,
	idPrefix: 'worker',
})

const emit = defineEmits<{
	close: []
	save: [worker: Worker]
	update: [worker: Worker]
}>()

const { editUrl, editPollInterval, headerRows, canSave, resetForm, addHeader, removeHeader, toWorker } = useWorkerForm(
	() => props.worker ?? null,
)

const handleSave = () => {
	if (!canSave.value) return
	emit('save', toWorker())
}

defineExpose({
	resetForm: () => resetForm(null),
})
</script>
