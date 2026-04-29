import type { config } from '@/composables/types';

export const toggleRealtime = async (enabled: boolean) => {
	try {
		const config: config = useBackendUrl()
		const response = await fetch($`${config.httphost}/settings/realtime`, {
			method: "POST",
			headers: {
				"Authorization": config.secret,
				"Session-Token": getCookie("session-token"),
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ enabled }),
		});
		// console.log(response);
		// console.log(enabled);
		if (!response.ok) throw new Error("Failed to toggle realtime");
	} catch (error) {
		console.error(error);
	}
}
