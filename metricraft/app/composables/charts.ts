import { Chart } from "chart.js";

export const createTrafficCongestionTrends = (
	canvas: HTMLCanvasElement,
	data: Record<string, Map<string, number>>
): Chart => {
	try {
		if (!data) {
			return new Chart(canvas, {
				type: 'bar',
				data: { labels: [], datasets: [] }
			});
		}
		const mapped: Map<string, Map<string, number>>[] = Object.entries(data)[1]
		const labels = Array.from(mapped.keys())
		console.log(labels);
		const values = Array.from(mapped.values()).map((v: any) => {
			if (v instanceof Map) return Array.from(v.values())[0];
			return v;
		});
		return new Chart(canvas, {
			type: 'bar',
			data: {
				labels,
				datasets: [{
					label: 'Congestion',
					data: values,
					backgroundColor: 'rgba(59, 130, 246, 0.5)',
					borderColor: 'rgb(59, 130, 246)',
					borderWidth: 1
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				scales: {
					y: { beginAtZero: true }
				}
			}
		});
	} catch (e) {
		console.error(e);
		return new Chart(canvas, {
			type: 'bar',
			data: { labels: [], datasets: [] }
		});
	}
}
