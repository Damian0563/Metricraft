import type { signPayload, config } from '@/composables/types';
export const welcome = async (): Promise<boolean | null> => {
	const config: config = useBackendUrl()
	try {
		const response = await fetch(`${config.httphost}/welcome`, {
			headers: {
				"Authorization": config.secret,
				"Session-Token": getCookie("session-token"),
			},
			method: "GET",
		})
		const data = await response.json()
		if (data.err) {
			throw data.err
		}
		return data.exists
	} catch (e) {
		console.log(e)
		throw "Something went wrong, Check your internet connection and try again."
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











