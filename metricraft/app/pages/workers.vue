<template>
	<div>
		<DashboardNav />
		<Popup :message="errorMessage ?? ''" @close="errorMessage = null" />
		<div class="w-full px-8 py-2">
			<div class="relative flex items-center justify-center mb-6">
				<button @click="goBack"
					class="absolute left-0 text-white hover:text-[#00F376] transition-colors duration-200 flex items-center gap-2 cursor-pointer">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
					</svg>
					Back to Dashboard
				</button>
				<h1 class="text-3xl font-bold text-center" style="color: #00F376;">Metricraft Workers</h1>
			</div>
			<div class="max-w-7xl mx-auto grid gap-8 lg:grid-rows-3">
				<div class="lg:row-span-1">
					<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100 h-full">
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
				</div>
				<div class="lg:row-span-2">
					<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
						<h2 class="text-xl font-semibold text-gray-800 mb-2">Configure worker</h2>
						<p class="text-sm text-gray-500 mb-6">
							Enter the upstream URL the worker should monitor and how frequently it should check the endpoint.
						</p>
						<div class="space-y-6">
							<div>
								<label for="worker-url" class="block text-sm font-medium text-gray-700 mb-2">Upstream URL</label>
								<input id="worker-url" v-model="workerUrl" type="url" placeholder="https://api.example.com/health"
									class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors" />
								<p class="mt-2 text-xs text-gray-500">
									The full address of the service the worker proxies traffic to and records metrics for.
								</p>
							</div>

							<div class="flex flex-col gap-8 xl:flex-row xl:justify-between">
								<div class="shrink-0">
									<label for="poll-interval" class="block text-sm font-medium text-gray-700 mb-2">Poll interval
										(minutes)</label>
									<input id="poll-interval" v-model.number="pollInterval" type="number" min="5" step="1"
										placeholder="10" max="60"
										class="w-full max-w-xs px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors" />
									<p class="mt-2 text-xs text-gray-500 max-w-xs">
										How often the worker sends a health check request to keep metrics fresh between real user traffic.
									</p>
								</div>
								<div class="flex-1 min-w-0">
									<div class="flex items-center justify-between mb-2">
										<label class="block text-sm font-medium text-gray-700">Headers</label>
										<button type="button" @click="addHeader"
											class="text-sm font-medium text-[#00F376] hover:text-[#00D96A] transition-colors cursor-pointer">
											+ Add header
										</button>
									</div>

									<div v-if="headerRows.length === 0"
										class="rounded-lg border border-dashed border-gray-200 px-4 py-6 text-center">
										<p class="text-sm text-gray-500">No custom headers configured.</p>
										<button type="button" @click="addHeader"
											class="mt-2 text-sm font-medium text-[#00F376] hover:text-[#00D96A] transition-colors cursor-pointer">
											Add your first header
										</button>
									</div>

									<div v-else class="space-y-2 max-h-48 overflow-y-auto pr-1">
										<div v-for="(header, index) in headerRows" :key="index" class="flex items-center gap-2">
											<input v-model="header.key" type="text" placeholder="Authorization"
												class="w-2/5 min-w-0 px-3 py-2 rounded-lg border border-gray-200 text-gray-800 text-sm focus:outline-none focus:border-[#00F376] transition-colors" />
											<input v-model="header.value" type="text" placeholder="Bearer token…"
												class="flex-1 min-w-0 px-3 py-2 rounded-lg border border-gray-200 text-gray-800 text-sm focus:outline-none focus:border-[#00F376] transition-colors" />
											<button type="button" @click="removeHeader(index)"
												class="shrink-0 p-2 text-gray-400 hover:text-red-500 transition-colors cursor-pointer"
												aria-label="Remove header">
												<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
													<path fill-rule="evenodd"
														d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
														clip-rule="evenodd" />
												</svg>
											</button>
										</div>
									</div>

									<p class="mt-2 text-xs text-gray-500">
										Optional HTTP headers included on each poll request, for example auth tokens or custom
										<code class="text-gray-600">Accept</code> values.
									</p>
								</div>
							</div>
						</div>

						<div class="mt-10 pt-6 border-t border-gray-100 flex justify-end">
							<button @click="save" :disabled="!canSave"
								class="px-8 py-3 bg-[#00F376] text-gray-900 font-bold rounded-lg hover:bg-[#00D96A] transition-all shadow-lg disabled:opacity-50 disabled:cursor-not-allowed uppercase tracking-wider text-sm cursor-pointer">
								Save worker
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { getCookie } from '@/composables/helpers'
import { saveWorker } from '@/composables/workers'
import type { Worker } from '@/composables/types'
type HeaderRow = { key: string; value: string }
const workerUrl = ref('')
const pollInterval = ref<number | null>(10)
const headerRows = ref<HeaderRow[]>([])
const errorMessage = ref<string | null>(null)
const canSave = computed(() => {
	const url = workerUrl.value.trim()
	const interval = pollInterval.value
	return url.length > 0 && interval !== null && interval >= 5
})

const buildHeaders = (): Record<string, string> => {
	const headers: Record<string, string> = {}
	for (const { key, value } of headerRows.value) {
		const trimmedKey = key.trim()
		const trimmedValue = value.trim()
		if (trimmedKey && trimmedValue) {
			headers[trimmedKey] = trimmedValue
		}
	}
	return headers
}

const addHeader = () => {
	headerRows.value.push({ key: '', value: '' })
}

const removeHeader = (index: number) => {
	headerRows.value.splice(index, 1)
}

const save = async () => {
	if (!canSave.value) return
	const worker: Worker = {
		url: workerUrl.value,
		pollInterval: pollInterval.value!,
		headers: buildHeaders(),
	}
	try {
		await saveWorker(worker)
	} catch (_) {
		errorMessage.value = "Something went wrong while saving the worker. Please try again."
	}
}

const goBack = () => {
	navigateTo('/dashboard')
}

onMounted(() => {
	const cookie = getCookie('session-token')
	if (!cookie) {
		navigateTo('/')
	}
})
</script>
