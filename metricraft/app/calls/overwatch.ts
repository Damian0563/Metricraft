import type { CustomMetric } from '@/composables/types/additional'

export const addCustomMetric = async (metric: CustomMetric) => {
	const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
	return await useApi()(`/overwatch/metrics/add?timezone=${encodeURIComponent(timezone)}`, {
		method: "POST",
		body: metric,
	})
}


