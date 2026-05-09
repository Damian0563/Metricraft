import type { dashboardInitPayload, config } from '@/composables/types';
import { getCookie } from "@/composables/helpers";

export const getDashboard = async (): Promise<dashboardInitPayload> => {
	try {
		const config: config = useBackendUrl()
		const headers = {
			"Authorization": config.secret,
			"Session-Token": useCookie("session-token").value || getCookie("session-token") || "",
		}
		const data = await $fetch<dashboardInitPayload>(`${config.httphost}/dashboard/init`, {
			method: "GET",
			headers,
			responseType: "json",
		})
		return data
	} catch (e) {
		return {
			appName: "",
			signedSecret: "",
			error: "Something went wrong, Check your internet connection and try again.",
			settings: {
				realtime: false,
				retention: 30,
				enabled: {},
			},
		}
	}
}
