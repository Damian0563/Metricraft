<template>
	<div>
		<Notice :message="errorMessage ?? ''" @close="errorMessage = null" />
		<Popup :message="cleanMessage" @close="cleanMessage = ''" />
		<Spinner :loading="saving" />
		<WorkerEditor :open="showWorkerEditor" :worker="editingWorker" @close="closeWorkerEditor"
			@save="handleWorkerUpdate" />
		<div class="relative flex items-center justify-center mb-6">
			<h1 class="text-3xl font-bold text-center" style="color: #00F376;">Metricraft Workers</h1>
		</div>
		<div class="max-w-8xl mx-auto grid gap-4 lg:grid-cols-[minmax(0,1fr)_40rem] lg:items-stretch">
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
			<div
				class="bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden lg:sticky lg:top-8 lg:h-full lg:grid lg:grid-rows-2">
				<section class="min-h-0 flex flex-col border-b border-gray-100">
					<div class="shrink-0 px-6 py-5 border-b border-gray-100">
						<h2 class="text-xl font-semibold text-gray-800">Configured workers</h2>
						<p class="text-sm text-gray-500 mt-1">
							{{ workersList.length }} worker{{ workersList.length === 1 ? '' : 's' }} monitoring your endpoints.
						</p>
					</div>
					<div v-if="workersList.length > 0" class="min-h-0 flex-1 overflow-y-auto">
						<ClientOnly>
							<AnimatePresence>
								<motion.div v-for="(worker, index) in workersList" :key="worker.url" layout
									:initial="{ opacity: 0, height: 0 }" :animate="{ opacity: 1, height: 'auto' }"
									:exit="{ opacity: 0, height: 0, x: 32, transition: { duration: 0.22, ease: 'easeIn' } }"
									:transition="{ type: 'spring', stiffness: 480, damping: 34, delay: Math.min(index * 0.045, 0.27) }"
									class="flex items-center gap-3 px-6 py-4 border-b border-gray-100 last:border-b-0 hover:bg-gray-50 transition-colors">
									<p class="min-w-0 flex-1 font-medium text-gray-900 text-sm break-all leading-snug">{{
										truncateUrl(worker.url, 30) }}</p>
									<div class="flex shrink-0 items-center gap-2">
										<span v-if="newWorkerStatusCodes[worker.url] !== undefined"
											:title="`Status code returned by the server: ${newWorkerStatusCodes[worker.url]}`"
											:class="['inline-flex items-center gap-1 px-2.5 py-1 text-xs font-bold rounded-full border', statusCodeClass(newWorkerStatusCodes[worker.url]!)]">
											<svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20"
												fill="currentColor">
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
								</motion.div>
							</AnimatePresence>
						</ClientOnly>
					</div>
					<div v-else class="min-h-0 flex-1 px-6 py-10 text-center">
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
				</section>
				<section class="min-h-0 flex flex-col bg-gray-50/60">
					<div class="shrink-0 flex items-start gap-3 px-6 py-5 border-b border-gray-100">
						<div
							class="h-10 w-10 shrink-0 rounded-full bg-[#00F376]/10 flex items-center justify-center text-[#00B35C]">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"
								aria-hidden="true">
								<path d="M2.94 6.34A2 2 0 014.73 5h10.54a2 2 0 011.79 1.34L10 10.76 2.94 6.34z" />
								<path d="M18 8.08l-7.47 4.67a1 1 0 01-1.06 0L2 8.08V14a2 2 0 002 2h12a2 2 0 002-2V8.08z" />
							</svg>
						</div>
						<div class="min-w-0 flex-1">
							<div class="flex items-center justify-between gap-3">
								<h2 class="text-sm font-semibold text-gray-800">Downtime notifications</h2>
								<span class="shrink-0 text-xs font-medium text-gray-500">
									{{ selectedRecipients.length }}/{{ users.length }} selected
								</span>
							</div>
							<p class="text-xs text-gray-500 mt-1">
								Choose who receives an email when a worker returns a status outside 200–299.
							</p>
							<div v-if="users.length > 0" class="flex items-center gap-3 mt-3">
								<button type="button" @click="selectAllRecipients"
									class="text-xs font-semibold text-[#00A652] hover:text-[#008C45] cursor-pointer">
									Select all
								</button>
								<button type="button" @click="clearRecipients"
									class="text-xs font-semibold text-gray-500 hover:text-gray-700 cursor-pointer">
									Clear
								</button>
								<button type="button" @click="saveRecipients"
									class="ml-auto inline-flex items-center justify-center px-3 py-1.5 text-xs font-semibold rounded-md bg-[#00F376] text-gray-900 hover:bg-[#00D96A] transition-colors cursor-pointer shadow-sm">
									Save
								</button>
							</div>
						</div>
					</div>
					<div v-if="users.length > 0" class="min-h-0 flex-1 overflow-y-auto p-4 space-y-2">
						<ClientOnly>
							<AnimatePresence>
								<motion.label v-for="(user, index) in users" :key="user.mail" layout
									:initial="{ opacity: 0, height: 0 }" :animate="{ opacity: 1, height: 'auto' }"
									:exit="{ opacity: 0, height: 0, x: 32, transition: { duration: 0.22, ease: 'easeIn' } }"
									:transition="{ type: 'spring', stiffness: 480, damping: 34, delay: Math.min(index * 0.045, 0.27) }"
									class="flex items-center gap-3 rounded-lg border bg-white px-4 py-3 transition-colors cursor-pointer"
									:class="isRecipientSelected(user.mail) ? 'border-[#00F376] ring-1 ring-[#00F376]/30' : 'border-gray-200 hover:border-gray-300'">
									<input :checked="isRecipientSelected(user.mail)" type="checkbox"
										class="h-4 w-4 rounded border-gray-300 accent-[#00D96A]" @change="toggleRecipient(user.mail)" />
									<div
										class="h-8 w-8 shrink-0 rounded-full bg-[#00F376]/10 flex items-center justify-center text-xs font-semibold text-[#00A652]">
										{{ user.initials }}
									</div>
									<span class="min-w-0 flex-1 truncate text-sm font-medium text-gray-800">{{ user.mail }}</span>
								</motion.label>
							</AnimatePresence>
						</ClientOnly>
					</div>
					<p v-else class="px-6 py-8 text-center text-sm text-gray-400">
						No active notification recipients are available.
					</p>
				</section>
			</div>
		</div>
		<div class="w-full bg-white rounded-xl shadow-xl border border-gray-100 p-8 mt-4 flex flex-col gap-8">
			<div v-if="workerUptimes.length > 0">
				<h2 class="text-xl font-semibold text-gray-800">Worker uptime</h2>
				<WorkerUptimeGraph v-for="entry in workerUptimes" :key="entry.url" :url="entry.url" :data="entry.data" />
			</div>
			<p v-else class="text-sm text-gray-500">No workers configured yet.</p>
		</div>
	</div>
</template>

<script setup lang="ts">
definePageMeta({
	layout: 'dashboard',
})
import { parseApiError, truncateUrl } from '@/composables/helpers'
import { getTeamUsers } from "@/calls/invite";
import { motion, AnimatePresence } from "motion-v";
import { getExistingWorkers, saveWorker, updateWorker, deleteWorkerEntry, getWorkerUptime } from '@/composables/workers'
import type { Worker, WorkerUptimeData, TeamUser } from '@/composables/types/additional'
const errorMessage = ref<string | null>(null)
const cleanMessage = ref<string>('')
const { data: existingWorkers, error: fetchError } = await useAsyncData<Worker[]>('existingWorkers', () => getExistingWorkers(), { default: () => [] })
const workersList = ref<Worker[]>([])
const showWorkerEditor = ref(false)
const editingWorker = ref<Worker | null>(null)
const originalWorkerUrl = ref<string | null>(null)
const saving = ref(false)
const workerUptimes = ref<{ url: string, data: WorkerUptimeData }[]>([])
const newWorkerFormRef = ref<{ resetForm: () => void } | null>(null)
const newWorkerStatusCodes = ref<Record<string, number>>({})
const { data: teamUsers, error: teamUsersError } = await useAsyncData<TeamUser[]>('teamUsers', () => getTeamUsers(), { default: () => [] })
watch(teamUsersError, () => errorMessage.value = teamUsersError.value?.message ?? '')
const users = computed(() => teamUsers.value?.filter(user => user.status) ?? [])
const selectedRecipients = ref<string[]>([])
const recipientSelectionInitialized = ref(false)
watch(users, (availableUsers) => {
	const notificationEnabledMails = new Set(
		availableUsers
			.filter(user => user.receiveNotifications)
			.map(user => user.mail),
	)
	if (!recipientSelectionInitialized.value) {
		selectedRecipients.value = [...notificationEnabledMails]
		recipientSelectionInitialized.value = availableUsers.length > 0
		return
	}
	selectedRecipients.value = selectedRecipients.value.filter(mail => notificationEnabledMails.has(mail))
}, { immediate: true })

watch(existingWorkers, async (workers) => {
	if (workers) {
		workersList.value = workers
		const results = await Promise.allSettled(workers.map(async (worker) => {
			const res: WorkerUptimeData = await getWorkerUptime(worker.url, worker.pollInterval)
			return { url: worker.url, data: res }
		}))
		workerUptimes.value = results
			.filter((result): result is PromiseFulfilledResult<{ url: string, data: WorkerUptimeData }> => result.status === 'fulfilled')
			.map(result => result.value)
	}
}, { immediate: true })

if (fetchError.value) {
	errorMessage.value = 'Failed to load workers.'
}

const isRecipientSelected = (mail: string): boolean => selectedRecipients.value.includes(mail)
const toggleRecipient = (mail: string) => {
	selectedRecipients.value = isRecipientSelected(mail)
		? selectedRecipients.value.filter(recipient => recipient !== mail)
		: [...selectedRecipients.value, mail]
}
const selectAllRecipients = () => {
	selectedRecipients.value = users.value.map(user => user.mail)
}
const clearRecipients = () => {
	selectedRecipients.value = []
}
const saveRecipients = async () => {
	const emails: { mail: string, receiveNotifications: boolean }[] = users.value.map(user => ({
		mail: user.mail,
		receiveNotifications: selectedRecipients.value.includes(user.mail)
	}));
	if (emails.length > 0) {
		saving.value = true
		try {
			await saveNotificationRecipients(emails)
		} catch (e) {
			errorMessage.value = parseApiError(e, 'Something went wrong while saving the worker. Please try again.')
		} finally {
			saving.value = false
		}
	}
}
const statusCodeClass = (statusCode: number): string => {
	if (statusCode >= 200 && statusCode < 300) return 'bg-[#00F376]/10 text-[#00B35C] border-[#00F376]'
	if (statusCode >= 300 && statusCode < 400) return 'bg-blue-50 text-blue-600 border-blue-300'
	if (statusCode >= 400 && statusCode < 500) return 'bg-amber-50 text-amber-600 border-amber-300'
	return 'bg-red-50 text-red-600 border-red-300'
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
</script>
