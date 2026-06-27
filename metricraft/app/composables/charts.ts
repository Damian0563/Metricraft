import { Chart } from "chart.js";
import { type ChartData } from "~/composables/types";
import { ColorPicker } from "~/composables/colorpicker";
import { ChoroplethChart, topojson } from 'chartjs-chart-geo';
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

export const createP95Latency = (
	canvas: HTMLCanvasElement,
	data: any
): Chart => {
	try {
		return new Chart(canvas, {
			type: 'bar',
			data: { labels: [], datasets: [] }
		})
	} catch (e) {
		return new Chart(canvas, {
			type: 'bar',
			data: { labels: [], datasets: [] }
		})
	}
}

export const createGeographicalTraffic = async (
	canvas: HTMLCanvasElement,
	data: any
): Promise<ChoroplethChart> => {
	try {
		const world = await fetch('https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json').then((r) => r.json());
		const countries = (topojson.feature(world, world.objects.countries) as any).features
			.filter((feature: any) => feature.properties.name !== 'Antarctica');
		const mapped = new Map<string, number>(Object.entries(data.distribution.values));
		const points: any[] = countries.map((feature: any) => {
			const value = mapped.get(feature.properties.name) ?? 0;
			return { feature, value: value > 0 ? value : null };
		});
		const all: number = points.reduce((acc: number, p: any) => acc + (p.value ?? 0), 0);
		canvas.style.backgroundColor = '#0b0f17';
		canvas.style.borderRadius = '3%';
		return new ChoroplethChart(canvas, {
			data: {
				labels: countries.map((feature: any) => feature.properties.name),
				datasets: [{
					label: 'Traffic',
					outline: countries,
					data: points,
					borderColor: 'rgba(11,15,23,0.9)',
					borderWidth: 0.5,
					borderJoinStyle: 'round',
					hoverBorderColor: '#00F376',
					hoverBorderWidth: 1.5,
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				devicePixelRatio: Math.ceil(window.devicePixelRatio || 1),
				animation: { duration: 250 },
				layout: { padding: { top: 8, right: 8, bottom: 8, left: 8 } },
				showOutline: true,
				showGraticule: false,
				elements: {
					geoFeature: {
						outlineBorderColor: 'rgba(255,255,255,0.08)',
						outlineBorderWidth: 0.75,
					},
				},
				scales: {
					projection: { axis: 'x', projection: 'mercator' },
					color: {
						axis: 'x',
						quantize: 6,
						missing: '#2a2f3a',
						interpolate: (v: number) => {
							const t = Math.max(0, Math.min(1, v));
							const g = Math.round(70 + 173 * t);
							const b = Math.round(34 + 84 * t);
							return `rgb(0, ${g}, ${b})`;
						},
						legend: {
							position: 'bottom-right',
							align: 'right',
							length: 180,
							width: 12,
							margin: 12,
							indicatorWidth: 12,
						},
						ticks: {
							color: '#e5e7eb',
							font: { weight: 'bold', size: 11 },
						},
					},
				},
				plugins: {
					legend: { display: false },
					tooltip: {
						enabled: true,
						displayColors: false,
						padding: 10,
						titleFont: { weight: 'bold', size: 13 },
						bodyFont: { size: 12 },
						callbacks: {
							title: (items) => (items[0]?.raw as any)?.feature?.properties?.name ?? '',
							label: (item) => {
								const value = (item.raw as any)?.value ?? 0;
								const share = all > 0 ? ((value / all) * 100).toFixed(1) : '0.0';
								return [`Traffic: ${value.toLocaleString()}`, `Share of total: ${share}%`];
							},
						},
					},
				},
			},
		})
	} catch (e) {
		console.error(e)
		return new ChoroplethChart(canvas, {
			data: { labels: [], datasets: [] }
		})
	}

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
