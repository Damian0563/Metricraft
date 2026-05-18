import { Chart } from "chart.js";

export const createTrafficCongestionTrends = (
	canvas: HTMLCanvasElement,
	data: Record<string, Map<string, number>>
): Chart => {
	try {
		if (!data) throw new Error('Data is empty');
		const labels: string[] | undefined = Object.entries(data).map(([_, v]) => Object.keys(v))[0];
		if (!labels) throw new Error('Labels are empty');
		const values = []
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
