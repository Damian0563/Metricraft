import type { config } from '@/composables/types';

export const toggleRealtime = async (enabled: boolean) => {
	try {
		const config: config = useBackendUrl()
		const headers = {
			"Authorization": config.secret,
			"Session-Token": getCookie("session-token"),
			"Content-Type": "application/json",
		}
		await $fetch(`${config.httphost}/settings/realtime`, {
			method: "POST",
			headers,
			body: { enabled },
		});
	} catch (error) {
		console.error(error);
	}
}

export const changeRetention = async (retention: number) => {
	try {
		const config: config = useBackendUrl()
		const headers = {
			"Authorization": config.secret,
			"Session-Token": getCookie("session-token"),
			"Content-Type": "application/json",
		}
		await $fetch(`${config.httphost}/settings/retention`, {
			method: "POST",
			headers,
			body: { retention },
		});
	} catch (error) {
		console.error(error);
	}
}

export const changeDerivedMetrics = async (metrics: { name: string; enabled: boolean }[]) => {
	try {
		const config: config = useBackendUrl()
		const headers = {
			"Authorization": config.secret,
			"Session-Token": getCookie("session-token"),
			"Content-Type": "application/json",
		}
		await $fetch<Response>(`${config.httphost}/settings/metrics`, {
			method: "POST",
			headers,
			body: metrics,
		});
	} catch (error) {
		console.error(error);
	}
}
