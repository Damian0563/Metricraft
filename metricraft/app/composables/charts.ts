import { Chart } from "chart.js";
import { type ChartData } from "~/composables/types";
import { ColorPicker } from "~/composables/colorpicker";
import { createAdditionalCongestionData } from "./chartUtils";

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

export const createGeographicalTraffic = (
	canvas: HTMLCanvasElement,
	data: Map<string, number>
): Chart => {
	return new Chart(canvas, {
		type: 'bar',
		data: { labels: [], datasets: [] }
	})
}

export const createTrafficCongestionTrends = (
	canvas: HTMLCanvasElement,
	data: TrafficCongestionData,
	colorPicker: ColorPicker
): ChartData => {
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
		let chart: Chart = new Chart(canvas, {
			type: 'bar',
			data: {
				labels,
				datasets
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				devicePixelRatio: Math.ceil(window.devicePixelRatio || 1),
				animation: { duration: 250 },
				interaction: { mode: 'index', intersect: false },
				layout: { padding: { right: 8, left: 2, bottom: 2 } },
				scales: {
					x: {
						stacked: true,
						grid: { display: false },
						border: { display: false },
						ticks: { autoSkip: false, color: 'black', maxRotation: 90, padding: 6, font: { weight: 'bold' } },
						title: { display: true, text: 'Time', color: 'black', font: { weight: 'bold', size: 18 }, padding: { top: 2 } },
					},
					y: {
						stacked: true,
						beginAtZero: true,
						border: { display: false },
						grid: { color: 'rgba(0,0,0,0.06)', drawTicks: false },
						ticks: { padding: 8, precision: 0, font: { weight: 'bold' } },
						title: { display: true, text: 'Requests', color: 'black', font: { weight: 'bold', size: 18 } },
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
						filter: (items) => items.raw !== 0,
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
		return { chart, additionalData: createAdditionalCongestionData(urlDataMap, colorPicker) };
	} catch (e) {
		return {
			chart: new Chart(canvas, {
				type: 'bar',
				data: { labels: [], datasets: [] }
			}), additionalData: null
		}
	}
}
