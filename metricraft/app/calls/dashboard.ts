import type { dashboardInitPayload } from '@/composables/types';

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


export const fetchMetric = async (metric: string, timeframe: string, persist: boolean = false): Promise<Map<string, number>> => {
	return await useApi()<Map<string, number>>(`/dashboard/fetch?persist=${persist}`, {
		method: "GET",
		query: { metric, timeframe },
	})
}
