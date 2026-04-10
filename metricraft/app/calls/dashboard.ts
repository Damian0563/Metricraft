import type { dashboardInitPayload } from '@/composables/types';

export const getDashboard = async (cookie: string): Promise<dashboardInitPayload> => {
	const SECRET = useBackendUrl()
	const response = await fetch(`http://localhost:8080/dashboard/init`, {
		headers: {
			"Authorization": SECRET,
			"Session-Token": cookie,
		},
		method: "GET",
	})
	const data = await response.json()
	return data.payload || data
}
