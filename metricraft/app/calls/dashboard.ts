import type { dashboardInitPayload, config } from '@/composables/types';

export const getDashboard = async (cookie: string): Promise<dashboardInitPayload> => {
	const config: config = useBackendUrl()
	const response = await fetch(`${config.httphost}/dashboard/init`, {
		headers: {
			"Authorization": config.secret,
			"Session-Token": cookie,
		},
		method: "GET",
	})
	let data = await response.json()
	if (!response.ok) {
		data.error = response.statusText
	}
	console.log(data)
	return data.payload || data
}
