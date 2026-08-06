import type { dashboardInitPayload } from '@/composables/types/views'
import type { CustomMetricResponse, MetricData } from '@/composables/types/additional'

export const getDashboard = async (): Promise<dashboardInitPayload> => {
	try {
		return await useApi()<dashboardInitPayload>("/dashboard/init", {
			method: "GET",
			responseType: "json",
		})
	} catch (e) {
		return {
			appName: "",
			signedSecret: "",
			error: "Something went wrong, Check your internet connection and try again.",
			urls: [],
			settings: {
				realtime: false,
				retention: 30,
				enabled: {},
			},
		}
	}
}


export const getUrls = async (): Promise<string[]> => {
	try {
		return await useApi()<string[]>("/dashboard/urls", {
			method: "GET",
			responseType: "json",
		})
	} catch (e) {
		return []
	}
}


export const fetchMetric = async (metric: string, timeframe: string, errorMessage: Ref<string>, persist: boolean = false): Promise<Map<string, number>> => {
	const userZone = Intl.DateTimeFormat().resolvedOptions().timeZone
	try {
		return await useApi()<Map<string, number>>(`/dashboard/fetch?persist=${persist}&timezone=${encodeURIComponent(userZone)}`, {
			method: "GET",
			query: { metric, timeframe },
		})
	} catch (e) {
		errorMessage.value = `Something went wrong while fetching "${metric}" metrics.`
		return new Map<string, number>()
	}
}

export const fetchCustomMetrics = async (errorMessage: Ref<string>): Promise<MetricData[]> => {
	const userZone = Intl.DateTimeFormat().resolvedOptions().timeZone
	try {
		const response = await useApi()<CustomMetricResponse>(`/dashboard/custom/fetch?timezone=${encodeURIComponent(userZone)}`, {
			method: "GET",
			responseType: "json",
		})
		if (response.errors?.length) {
			errorMessage.value = response.errors.join('')
		}
		return response.metrics ?? []
	} catch (_) {
		errorMessage.value = `Something went wrong while fetching custom metrics.`
		return []
	}
}
