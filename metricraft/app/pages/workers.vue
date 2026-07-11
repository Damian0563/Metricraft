<template>
	<div class="min-h-screen pl-20">
		<DashboardNav />
		<Notice :message="errorMessage ?? ''" @close="errorMessage = null" />
		<Popup :message="cleanMessage" @close="cleanMessage = ''" />
		<Spinner :loading="saving" />
		<WorkerEditor :open="showWorkerEditor" :worker="editingWorker" @close="closeWorkerEditor"
			@save="handleWorkerUpdate" />
		<div class="min-w-0 px-8 py-2">
			<div class="relative flex items-center justify-center mb-6">
				<h1 class="text-3xl font-bold text-center" style="color: #00F376;">Metricraft Workers</h1>
			</div>
			<div class="max-w-8xl mx-auto grid gap-4 lg:grid-cols-[minmax(0,1fr)_40rem] lg:items-start">
				<div class="flex flex-col gap-4 min-w-0">
					<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
						<h2 class="text-xl font-semibold text-gray-800 mb-4">What are workers?</h2>
						<p class="text-sm text-gray-600 leading-relaxed mb-4">
							Workers are lightweight reverse proxies that sit in front of your application and capture HTTP
							traffic — request method, URL, status codes, latency, and headers — without changing how your app
							runs.
						</p>
						<p class="text-sm text-gray-600 leading-relaxed">
							Each worker forwards requests to your upstream service, streams metrics to the Metricraft backend,
							and keeps your dashboards updated with real observability data.
						</p>
						<p class="text-sm text-gray-500">
							Point a worker at an endpoint and set how often it should poll so uptime and performance metrics stay
							current even when traffic is low.
						</p>
					</div>
					<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
						<h2 class="text-xl font-semibold text-gray-800 mb-2">Configure worker</h2>
						<p class="text-sm text-gray-500 mb-6">
							Enter the upstream URL the worker should monitor and how frequently it should check the endpoint.
						</p>
						<WorkerForm ref="newWorkerFormRef" :saving="saving" show-helper-text save-label="Save worker"
							@save="handleNewWorkerSave" />
					</div>
				</div>
				<div class="bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden lg:sticky lg:top-8 h-full">
					<div class="px-6 py-5 border-b border-gray-100">
						<h2 class="text-xl font-semibold text-gray-800">Configured workers</h2>
						<p class="text-sm text-gray-500 mt-1">
							{{ workersList.length }} worker{{ workersList.length === 1 ? '' : 's' }} monitoring your endpoints.
						</p>
					</div>
					<div v-if="workersList.length > 0" class="max-h-[calc(100vh-12rem)] overflow-y-auto">
						<div v-for="worker in workersList" :key="worker.url"
							class="flex items-center gap-3 px-6 py-4 border-b border-gray-100 last:border-b-0 hover:bg-gray-50 transition-colors">
							<p class="min-w-0 flex-1 font-medium text-gray-900 text-sm break-all leading-snug">{{ worker.url }}</p>
							<div class="flex shrink-0 items-center gap-2">
								<span v-if="newWorkerStatusCodes[worker.url] !== undefined"
									:title="`Status code returned by the server: ${newWorkerStatusCodes[worker.url]}`"
									:class="['inline-flex items-center gap-1 px-2.5 py-1 text-xs font-bold rounded-full border', statusCodeClass(newWorkerStatusCodes[worker.url]!)]">
									<svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
										<path fill-rule="evenodd"
											d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V7z"
											clip-rule="evenodd" />
									</svg>
									{{ newWorkerStatusCodes[worker.url] }}
								</span>
								<button type="button" @click="openWorkerEditor(worker)"
									class="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-semibold rounded-lg border border-[#00F376] text-[#00B35C] bg-[#00F376]/10 hover:bg-[#00F376]/20 transition-colors cursor-pointer">
									<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
										<path
											d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
									</svg>
									Edit worker
								</button>
								<button type="button" @click="deleteWorker(worker.url)"
									class="px-4 py-2 text-sm font-semibold rounded-lg bg-red-50 text-red-600 hover:bg-red-100 transition-colors cursor-pointer">
									Delete
								</button>
							</div>
						</div>
					</div>
					<div v-else class="px-6 py-10 text-center">
						<div class="mx-auto mb-3 h-10 w-10 rounded-full bg-gray-100 flex items-center justify-center text-gray-400">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd"
									d="M11.3 1.046A1 1 0 0112 2v5h4a1 1 0 01.82 1.573l-7 10A1 1 0 018 18v-5H4a1 1 0 01-.82-1.573l7-10a1 1 0 011.12-.38z"
									clip-rule="evenodd" />
							</svg>
						</div>
						<p class="text-sm text-gray-500">No workers configured yet.</p>
						<p class="text-xs text-gray-400 mt-1">Saved workers will appear here.</p>
					</div>
				</div>
			</div>
			<div v-if="workerUptimes.length > 0"
				class="w-full bg-white rounded-xl shadow-xl border border-gray-100 p-8 mt-4 flex flex-col gap-8">
				<h2 class="text-xl font-semibold text-gray-800">Worker uptime</h2>
				<WorkerUptimeGraph v-for="entry in workerUptimes" :key="entry.url" :url="entry.url"
					:data="entry.data" />
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { getCookie, parseApiError } from '@/composables/helpers'
import { getExistingWorkers, saveWorker, updateWorker, deleteWorkerEntry, getWorkerUptime } from '@/composables/workers'
import type { Worker, WorkerUptimeData } from '@/composables/types'
const errorMessage = ref<string | null>(null)
const cleanMessage = ref<string>('')
const { data: existingWorkers, error: fetchError } = await useAsyncData<Worker[]>('existingWorkers', () => getExistingWorkers())
const workersList = ref<Worker[]>([])
const showWorkerEditor = ref(false)
const editingWorker = ref<Worker | null>(null)
const originalWorkerUrl = ref<string | null>(null)
const saving = ref(false)
const workerUptimes = ref<{ url: string, data: WorkerUptimeData }[]>([])
const newWorkerFormRef = ref<{ resetForm: () => void } | null>(null)
const newWorkerStatusCodes = ref<Record<string, number>>({})
const statusCodeClass = (statusCode: number): string => {
	if (statusCode >= 200 && statusCode < 300) return 'bg-[#00F376]/10 text-[#00B35C] border-[#00F376]'
	if (statusCode >= 300 && statusCode < 400) return 'bg-blue-50 text-blue-600 border-blue-300'
	if (statusCode >= 400 && statusCode < 500) return 'bg-amber-50 text-amber-600 border-amber-300'
	return 'bg-red-50 text-red-600 border-red-300'
}
watch(existingWorkers, async (workers) => {
	if (workers) {
		workersList.value = workers
		workerUptimes.value = await Promise.all(workers.map(async (worker) => {
			const res: WorkerUptimeData = await getWorkerUptime(worker.url)
			return { url: worker.url, data: res }
		}))
	}
}, { immediate: true })

if (fetchError.value) {
	errorMessage.value = 'Failed to load workers.'
}

const openWorkerEditor = (worker: Worker) => {
	editingWorker.value = worker
	originalWorkerUrl.value = worker.url
	showWorkerEditor.value = true
}

const closeWorkerEditor = () => {
	showWorkerEditor.value = false
	editingWorker.value = null
	originalWorkerUrl.value = null
}

const handleWorkerUpdate = async (updatedWorker: Worker) => {
	const originalUrl = originalWorkerUrl.value
	if (!originalUrl) return
	const index = workersList.value.findIndex((worker) => worker.url === originalUrl)
	if (index === -1) return
	saving.value = true
	try {
		const res: { success: boolean; err: string; statusCode: number } = await updateWorker(updatedWorker)
		if (!res.success) {
			errorMessage.value = res.err
		}
		workersList.value = [
			...workersList.value.slice(0, index),
			updatedWorker,
			...workersList.value.slice(index + 1),
		]
		if (newWorkerStatusCodes.value.hasOwnProperty(originalUrl) && newWorkerStatusCodes.value[originalUrl] !== undefined) {
			delete newWorkerStatusCodes.value[originalUrl]
		}
		newWorkerStatusCodes.value[updatedWorker.url] = res.statusCode
		cleanMessage.value = 'Worker updated successfully.'
	} catch (e: unknown) {
		errorMessage.value = parseApiError(e, 'Something went wrong while saving the worker. Please try again.')
	} finally {
		saving.value = false
		closeWorkerEditor()
	}

}

const handleNewWorkerSave = async (worker: Worker) => {
	saving.value = true
	try {
		const res: { success: boolean; err: string; statusCode: number } = await saveWorker(worker)
		if (!res.success) {
			errorMessage.value = res.err
		}
		workersList.value.push(worker)
		newWorkerStatusCodes.value[worker.url] = res.statusCode
		newWorkerFormRef.value?.resetForm()
	} catch (e: unknown) {
		errorMessage.value = parseApiError(e, 'Something went wrong while saving the worker. Please try again.')
	} finally {
		saving.value = false
	}
}

const deleteWorker = async (url: string) => {
	closeWorkerEditor()
	saving.value = true
	try {
		const res: { success: boolean; err: string } = await deleteWorkerEntry(url)
		if (!res.success) {
			errorMessage.value = res.err
			return
		}
		workersList.value = workersList.value.filter((worker) => worker.url !== url)
		delete newWorkerStatusCodes.value[url]
		if (originalWorkerUrl.value === url) {
			closeWorkerEditor()
		}
		cleanMessage.value = 'Worker deleted successfully.'
	} catch (e: unknown) {
		errorMessage.value = parseApiError(e, 'Something went wrong while deleting the worker. Please try again.')
	} finally {
		saving.value = false
	}
}

onMounted(() => {
	const cookie = getCookie('session-token')
	if (!cookie) {
		navigateTo('/')
	}
})
</script>
