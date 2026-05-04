import type { config } from '@/composables/types';

export const toggleRealtime = async (enabled: boolean) => {
	try {
		const config: config = useBackendUrl()
		const response = await fetch(`${config.httphost}/settings/realtime`, {
			method: "POST",
			headers: {
				"Authorization": config.secret,
				"Session-Token": getCookie("session-token"),
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ enabled }),
		});
		if (!response.ok) throw new Error("Failed to toggle realtime");
	} catch (error) {
		console.error(error);
	}
}

export const changeRetention = async (retention: number) => {
	try {
		const config: config = useBackendUrl()
		const response = await fetch(`${config.httphost}/settings/retention`, {
			method: "POST",
			headers: {
				"Authorization": config.secret,
				"Session-Token": getCookie("session-token"),
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ retention }),
		});
		console.log(response)
		if (!response.ok) throw new Error("Failed to change retention");
	} catch (error) {
		console.error(error);
	}
}

export const changeDerivedMetrics = async (metrics: { name: string; enabled: boolean }[]) => {
	try {
		const config: config = useBackendUrl()
		const response = await fetch(`${config.httphost}/settings/metrics`, {
			method: "POST",
			headers: {
				"Authorization": config.secret,
				"Session-Token": getCookie("session-token"),
				"Content-Type": "application/json",
			},
			body: { "metrics": JSON.stringify(metrics) },
		});
		console.log(response, metrics)
		if (!response.ok) throw new Error("Failed to change metrics");
	} catch (error) {
		console.error(error);
	}
}
