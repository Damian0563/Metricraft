import type { CustomMetric } from '@/composables/types/additional'

export const addCustomMetric = async (metric: Omit<CustomMetric, 'lastUpdate'>): Promise<void> => {
	const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
	return await useApi()(`/overwatch/metrics/add?timezone=${encodeURIComponent(timezone)}`, {
		method: "POST",
		body: metric,
	})
}

export const getCustomMetrics = async (): Promise<CustomMetric[]> => {
	const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
	return await useApi()(`/overwatch/metrics?timezone=${encodeURIComponent(timezone)}`)
}
