import type { signPayload } from '@/composables/types';
export const welcome = async (): Promise<boolean | null> => {
	const [SECRET, PORT] = useBackendUrl()
	const response = await fetch(`http://localhost:${PORT}/`, {
		headers: {
			"Authorization": SECRET,
			"Session-Token": getCookie("session-token"),
		},
		method: "GET",
	})
	const data = await response.json()
	if (data.err) {
		throw data.err
	}
	return data.exists
}

export const sign = async (payload: signPayload): Promise<boolean | null> => {
	const [SECRET, PORT] = useBackendUrl()
	const response = await fetch(`http://localhost:${PORT}/sign`, {
		headers: {
			"Authorization": SECRET,
		},
		method: "POST",
		body: JSON.stringify(payload),
	})
	const data = await response.json()
	if (data.err) {
		throw data.err
	}
	return data.token
}











