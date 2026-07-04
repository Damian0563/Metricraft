import type { Worker } from '@/composables/types'

export const saveWorker = async (worker: Worker): Promise<void> => {
	await useApi()("/dashboard/worker/new", {
		method: "POST",
		body: worker,
	})
}
