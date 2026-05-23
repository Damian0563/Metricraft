import { Chart } from "chart.js";

type StringInt32Map = {
	values?: Record<string, number>;
};
type CongestionEntry = {
	timerange: string;
	pairing?: StringInt32Map;
};
export type TrafficCongestionData = {
	values: CongestionEntry[];
};

export const createTrafficCongestionTrends = (
	canvas: HTMLCanvasElement,
	data: TrafficCongestionData,
	colorPicker: ColorPicker
): Chart => {
	const getUrlCounts = (pairing: StringInt32Map | undefined): Record<string, number> =>
		pairing?.values ?? {};
	try {
		if (!data?.values?.length) throw new Error('Data is empty');
		const labels: string[] = [];
		const urlDataMap = new Map<string, number[]>();
		const totalPoints = data.values.length;
		data.values.forEach((entry: CongestionEntry, pointIndex: number) => {
			labels.push(entry.timerange);
			for (const [url, count] of Object.entries(getUrlCounts(entry.pairing))) {
				if (!urlDataMap.has(url)) {
					urlDataMap.set(url, new Array<number>(totalPoints).fill(0));
				}
				urlDataMap.get(url)![pointIndex] = count;
			}
		});
		const datasets = Array.from(urlDataMap.entries()).map(([url, dataArray]) => {
			const color = colorPicker.getColorForUrl(url);
			return {
				label: url,
				data: dataArray,
				backgroundColor: color,
				hoverBackgroundColor: color,
				borderWidth: 0,
				borderRadius: 4,
				borderSkipped: false,
				categoryPercentage: 0.7,
				barPercentage: 0.9,
			};
		});
		return new Chart(canvas, {
			type: 'bar',
			data: {
				labels,
				datasets
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				animation: { duration: 250 },
				interaction: { mode: 'index', intersect: false },
				layout: { padding: { right: 8, left: 2, bottom: 4 } },
				scales: {
					x: {
						stacked: true,
						grid: { display: false },
						border: { display: false },
						ticks: { autoSkip: false, maxRotation: 90, padding: 6, font: { weight: 'bold' } },
						title: { display: true, text: 'Time', font: { weight: 'bold', size: 18 }, padding: { top: 8 } },
					},
					y: {
						stacked: true,
						beginAtZero: true,
						border: { display: false },
						grid: { color: 'rgba(0,0,0,0.06)', drawTicks: false },
						ticks: { padding: 8, precision: 0, font: { weight: 'bold' } },
						title: { display: true, text: 'Requests', font: { weight: 'bold', size: 18 } },
					}
				},
				plugins: {
					legend: {
						display: datasets.length <= 12,
						position: 'top',
						align: 'start',
						labels: {
							usePointStyle: true,
							pointStyle: 'circle',
							boxWidth: 8,
							boxHeight: 8,
							padding: 14,
						},
					},
					tooltip: {
						enabled: true,
						mode: 'index',
						intersect: false,
						itemSort: (a, b) => (b.parsed.y as number) - (a.parsed.y as number),
						callbacks: {
							footer: (items) => {
								const total = items.reduce((s, i) => s + (i.parsed.y as number), 0);
								return `Total: ${total}`;
							},
						},
					}
				},
			}
		});
	} catch (e) {
		return new Chart(canvas, {
			type: 'bar',
			data: { labels: [], datasets: [] }
		});
	}
}
