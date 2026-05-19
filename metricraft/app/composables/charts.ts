import { Chart } from "chart.js";

export const createTrafficCongestionTrends = (
	canvas: HTMLCanvasElement,
	data: Record<string, Map<string, number>>
): Chart => {
	try {
		if (!data) throw new Error('Data is empty');
		let labels: string[] = [];
		const map: Map<string, Map<string, number>> = new Map(
			Object.entries(data.values).map(([dateKey, innerWrapper]: [string, any]) => {
				const actualUrlCounts = innerWrapper.values ? innerWrapper.values : innerWrapper;
				return [
					dateKey,
					new Map(Object.entries(actualUrlCounts))
				];
			})
		);
		const urlDataMap = new Map<string, number[]>();
		const totalPoints = map.size;
		let pointIndex = 0;
		map.forEach((innerMap: Map<string, number>, key: string) => {
			labels.push(key);
			innerMap.forEach((count: number, url: string) => {
				if (!urlDataMap.has(url)) {
					urlDataMap.set(url, new Array(totalPoints).fill(0));
				}
				urlDataMap.get(url)![pointIndex] = count;
			});
			pointIndex++;
		});
		const datasets = Array.from(urlDataMap.entries()).map(([url, dataArray], _) => {
			return {
				label: url,
				data: dataArray,
				borderWidth: 1,
				backgroundColor: 'rgba(255, 99, 132, 0.5)',
			};
		});
		return new Chart(canvas, {
			type: 'bar',
			data: {
				labels,
				datasets
			},
			options: {
				interaction: {
					mode: 'nearest',
					intersect: true
				},
				responsive: true,
				maintainAspectRatio: false,
				scales: {
					x: { stacked: true },
					y: { beginAtZero: true, stacked: true }
				},
				plugins: {
					legend: {
						display: true,
						position: 'top'
					},
					tooltip: {
						enabled: true,
					}
				},
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
