import type { signPayload, config, welcomeResponse } from '@/composables/types';
import { getCookie } from "@/composables/helpers";
export const welcome = async (): Promise<boolean> => {
	try {
		const config: config = useBackendUrl()
		const headers = {
			"Authorization": config.secret,
			"Session-Token": useCookie("session-token").value || getCookie("session-token") || "",
		}
		const response = await $fetch<welcomeResponse>(`${config.httphost}/welcome`, {
			method: "GET", headers
		})
		if (response.err) throw (response.err)
		return Boolean(response.exists)
	} catch (e) {
		return false
	}
}

export const sign = async (payload: signPayload): Promise<boolean | null> => {
	const config: config = useBackendUrl()
	try {
		const response = await fetch(`${config.httphost}/sign`, {
			method: 'POST',
			headers: {
				"Authorization": config.secret,
			},
			body: JSON.stringify(payload),
		})
		const data = await response.json()
		if (data.err) {
			throw data.err
		}
		return data.token
	} catch (e) {
		console.log(e)
		throw "Something went wrong, Check your internet connection and try again."
	}
}











