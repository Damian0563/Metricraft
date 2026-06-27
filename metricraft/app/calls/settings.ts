export const toggleRealtime = async (enabled: boolean) => {
	try {
		await useApi()(`/settings/realtime`, {
			method: "POST",
			body: { enabled },
		});
	} catch (error) {
		console.error(error);
	}
}

export const changeRetention = async (retention: number) => {
	try {
		await useApi()(`/settings/retention`, {
			method: "POST",
			body: { retention },
		});
	} catch (error) {
		console.error(error);
	}
}

export const changeDerivedMetrics = async (metrics: { name: string; enabled: boolean; timeframe?: string }[]) => {
	try {
		await useApi()<Response>(`/settings/metrics`, {
			method: "POST",
			body: metrics,
		});
	} catch (error) {
		console.error(error);
	}
}
