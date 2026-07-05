import type { Worker } from '@/composables/types'

export const saveWorker = async (worker: Worker): Promise<{ success: boolean; err: string }> => {
	return await useApi()("/dashboard/worker/new", {
		method: "POST",
		body: worker,
	})
}


export const getExistingWorkers = async (): Promise<Worker[]> => {
	return await useApi()("/dashboard/worker/list", {
		method: "GET",
	})
}
