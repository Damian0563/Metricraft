import type { CustomMetric } from '@/composables/types/additional'

export const addCustomMetric = async (metric: Omit<CustomMetric, 'lastUpdate'>): Promise<void> => {
	return await useApi()('/overwatch/metrics/add', {
		method: "POST",
		body: metric,
	})
}

export const getCustomMetrics = async (): Promise<CustomMetric[]> => {
	return await useApi()('/overwatch/metrics')
}

export const deleteCustomMetric = async (metric: CustomMetric): Promise<void> => {
	return await useApi()(`/overwatch/metrics/delete`, {
		method: "POST",
		body: metric
	})
}


export const updateCustomMetric = async (original: CustomMetric, updated: CustomMetric): Promise<void> => {
	return await useApi()('/overwatch/metrics/update', {
		method: "PATCH",
		body: { original: original, updated: updated },
	})
}
