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
