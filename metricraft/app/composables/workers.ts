import type { Worker, HeaderRow } from '@/composables/types'

export const headersToRows = (headers?: Record<string, string>): HeaderRow[] => {
	if (!headers) return []
	return Object.entries(headers).map(([key, value]) => ({ key, value }))
}

export const buildHeadersFromRows = (headerRows: HeaderRow[]): Record<string, string> => {
	const headers: Record<string, string> = {}
	for (const { key, value } of headerRows) {
		const trimmedKey = key.trim()
		const trimmedValue = value.trim()
		if (trimmedKey && trimmedValue) {
			headers[trimmedKey] = trimmedValue
		}
	}
	return headers
}

export const isWorkerFormValid = (url: string, pollInterval: number | null): boolean => {
	return url.trim().length > 0 && pollInterval !== null && pollInterval >= 5 && pollInterval <= 60
}

export const workerFromForm = (url: string, pollInterval: number | null, headerRows: HeaderRow[]): Worker => ({
	url: url.trim(),
	pollInterval: pollInterval!,
	headers: buildHeadersFromRows(headerRows),
})

export const useWorkerForm = (source: () => Worker | null) => {
	const editUrl = ref('')
	const editPollInterval = ref<number | null>(10)
	const headerRows = ref<HeaderRow[]>([])
	const canSave = computed(() => isWorkerFormValid(editUrl.value, editPollInterval.value))
	const resetForm = (worker: Worker | null) => {
		if (!worker) {
			editUrl.value = ''
			editPollInterval.value = 10
			headerRows.value = []
			return
		}
		editUrl.value = worker.url
		editPollInterval.value = worker.pollInterval
		headerRows.value = headersToRows(worker.headers)
	}

	const addHeader = () => {
		headerRows.value.push({ key: '', value: '' })
	}

	const removeHeader = (index: number) => {
		headerRows.value.splice(index, 1)
	}
	const toWorker = (): Worker => workerFromForm(editUrl.value, editPollInterval.value, headerRows.value)
	watch(source, resetForm, { immediate: true })

	return {
		editUrl,
		editPollInterval,
		headerRows,
		canSave,
		resetForm,
		addHeader,
		removeHeader,
		toWorker,
	}
}

export const saveWorker = async (worker: Worker): Promise<{ success: boolean; err: string; statusCode: number }> => {
	return await useApi()("/dashboard/worker/new", {
		method: "POST",
		body: worker,
	})
}

export const updateWorker = async (worker: Worker): Promise<{ success: boolean; err: string; statusCode: number }> => {
	return await useApi()(`/dashboard/worker/update`, {
		method: "PATCH",
		body: worker,
	})
}

export const deleteWorkerEntry = async (workerUrl: string): Promise<{ success: boolean; err: string }> => {
	return await useApi()(`/dashboard/worker/delete`, {
		method: "DELETE",
		body: { url: workerUrl },
	})
}

export const getExistingWorkers = async (): Promise<Worker[]> => {
	return await useApi()("/dashboard/worker/list", {
		method: "GET",
	})
}
