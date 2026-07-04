import { Chart } from "chart.js";
import type { ChartData, TrafficCongestionData, ThroughputData, ThroughputEntry, DistributionData, additionalDataHeaders, CongestionEntry, StringInt32Map } from "~/composables/types";
import { ColorPicker } from "~/composables/colorpicker";
import { ChoroplethChart, topojson } from 'chartjs-chart-geo';
import { createAdditionalData } from "./chartUtils";

const emptyChart = (canvas: HTMLCanvasElement): Chart => {
	return new Chart(canvas, {
		type: 'bar',
		data: { labels: [], datasets: [] },
		options: { responsive: true, maintainAspectRatio: false },
	});
}
const formatTimeRangeTitle = (labels: string[], dataIndex: number, timeframe: string | null = null): string => {
	if (dataIndex < 0 || dataIndex >= labels.length) return '';
	if (timeframe !== "1d" && timeframe !== "7d") {
		const start = labels[dataIndex];
		return start ?? '';
	} else {
		const start = labels[dataIndex];
		const end = labels[dataIndex + 1];
		return end != null ? `${start} - ${end}` : `${start}-${labels[0]}`;
	}

};
const resolveHoveredIndex = (
	tooltipItems: Array<{ dataIndex?: number }>,
	chartInstance: Chart,
): number =>
	tooltipItems.find((item) => item.dataIndex != null)?.dataIndex ??
	chartInstance.getActiveElements()[0]?.index ??
	-1;

const emptyChoroplethChart = (canvas: HTMLCanvasElement): ChoroplethChart =>
	new ChoroplethChart(canvas, {
		data: { labels: [], datasets: [] },
	}) as ChoroplethChart;
const truncateUrl = (url: string, max = 20): string => {
	if (url.length <= max) return url;
	return `${url.slice(0, max - 1)}…`;
};

export const createThroughput = (
	canvas: HTMLCanvasElement,
	data: ThroughputData,
	timeframe: string,
): ChartData => {
	try {
		if (!data?.values?.length) throw new Error('Data is empty');
		const mapped = new Array<{ timerange: string, value: number }>();
		const labels: string[] = [];
		const values: number[] = [];
		data.values.forEach((entry: ThroughputEntry) => {
			const { timerange, value } = entry;
			const cleanedVal = value !== undefined ? value : 0;
			labels.push(timerange);
			values.push(cleanedVal);
			mapped.push({ timerange: timerange, value: cleanedVal });
		});
		let stepSize;
		switch (timeframe) {
			case "1d":
				stepSize = 6;
				break;
			case "7d":
				stepSize = 1;
				break;
			case "30d":
				stepSize = 2;
				break;
			default:
				stepSize = 3;
		}
		type ThroughputChart = Chart & { $hoveredIndex?: number | null };
		const chart: Chart = new Chart(canvas, {
			type: 'line',
			plugins: [
				{
					id: 'throughputLabelHover',
					afterInit(chart) {
						const throughputChart = chart as ThroughputChart;
						throughputChart.$hoveredIndex = null;
					},
					afterEvent(chart) {
						const throughputChart = chart as ThroughputChart;
						const active = chart.getActiveElements();
						const nextIndex = active[0]?.index ?? null;
						if (nextIndex === throughputChart.$hoveredIndex) return;
						throughputChart.$hoveredIndex = nextIndex;
						chart.update('none');
					},
				},
			],
			data: {
				labels,
				datasets: [{
					label: 'Throughput',
					data: values,
					showLine: true,
					tension: 0.3,
					borderColor: '#00C962',
					borderWidth: 2,
					backgroundColor: '#00F376',
					pointBackgroundColor: '#00F376',
					pointBorderColor: '#FFFFFF',
					pointBorderWidth: 2,
					pointRadius: 4,
					pointHoverRadius: 6,
					pointHoverBackgroundColor: '#00F376',
					pointHoverBorderColor: '#FFFFFF',
					pointHoverBorderWidth: 2,
				}],
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				devicePixelRatio: Math.ceil(window.devicePixelRatio || 1),
				animation: { duration: 250 },
				interaction: { mode: 'index', intersect: false },
				layout: { padding: { right: 8, left: 2, bottom: 2, top: 4 } },
				scales: {
					x: {
						grid: { display: false },
						border: { display: false },
						ticks: {
							autoSkip: false,
							maxRotation: 0,
							padding: 6,
							color: '#475569',
							font: { weight: 'bold', size: 11 },
							callback: (_value: string | number, index: number): string =>
								index % stepSize === 0 || index === 0 || (index === labels.length - 1 && (index - 1) % stepSize !== 0) ? (labels[index] ?? '') : '',
						},
						title: {
							display: true,
							text: 'Time',
							color: '#475569',
							font: { weight: 'bold', size: 13 },
							padding: { top: 2 },
						},
					},
					y: {
						beginAtZero: true,
						border: { display: false },
						grid: { color: 'rgba(0,0,0,0.06)', drawTicks: false },
						ticks: { padding: 8, precision: 0, font: { weight: 'bold' } },
						title: {
							display: true,
							text: 'Requests',
							color: '#475569',
							font: { weight: 'bold', size: 13 },
						},
					},
				},
				plugins: {
					legend: { display: false },
					tooltip: {
						enabled: true,
						displayColors: true,
						padding: 10,
						titleFont: { weight: 'bold', size: 13 },
						bodyFont: { size: 12 },
						callbacks: {
							title: function (tooltipItems) {
								return formatTimeRangeTitle(labels, resolveHoveredIndex(tooltipItems, this.chart), timeframe);
							},
							label: (item) => `Requests: ${Number(item.parsed.y).toLocaleString()}`,
							labelColor: () => ({
								borderColor: '#00F376',
								backgroundColor: '#00F376',
							}),
						},
					},
				},
			},
		});
		let addData: Array<{ timerange: string, value: number }> = new Array<{ timerange: string, value: number }>();
		if (timeframe === "1d") {
			values.forEach((value, index) => {
				addData.push({ timerange: formatTimeRangeTitle(labels, index), value: value });
			});
		} else {
			addData = mapped
		}
		const headers: additionalDataHeaders = { h1: 'Timerange', h2: 'Number of requests' };
		return { chart, additionalData: createAdditionalData(addData, headers) };
	} catch (e) {
		console.error(e)
		return { chart: emptyChart(canvas), additionalData: null };
	}
}

export const createUptimeScore = (
	canvas: HTMLCanvasElement,
	data: DistributionData,
	colorPicker: ColorPicker,
): ChartData => {
	const uptimeBarColor = (score: number): string => {
		if (score === 0) return '#DC2626';
		if (score >= 99.5) return '#00F376';
		if (score >= 95) return '#34D399';
		if (score >= 80) return '#FBBF24';
		if (score >= 50) return '#FB923C';
		return '#F87171';
	};
	try {
		const values = data?.distribution?.values;
		if (!values || !Object.keys(values).length) throw new Error('Data is empty');
		const entries: Array<[string, number]> = Object.entries(values).sort(([, a], [, b]) => a - b);
		const mapped: Map<string, number> = new Map(entries);
		const labels: string[] = Array.from(mapped.keys());
		const scores: number[] = Array.from(mapped.values());
		const isDown = scores.map((score) => score === 0);
		type UptimeChart = Chart & { $hoveredIndex?: number | null };
		const chart: Chart = new Chart(canvas, {
			type: 'bar',
			plugins: [
				{
					id: 'uptimeLabelHover',
					afterInit(chart) {
						const uptimeChart = chart as UptimeChart;
						uptimeChart.$hoveredIndex = null;
					},
					afterEvent(chart) {
						const uptimeChart = chart as UptimeChart;
						const active = chart.getActiveElements();
						const nextIndex = active[0]?.index ?? null;
						if (nextIndex === uptimeChart.$hoveredIndex) return;
						uptimeChart.$hoveredIndex = nextIndex;
						chart.update('none');
					},
				},
				{
					id: 'uptimeZeroMarker',
					afterDatasetsDraw(chart) {
						const { ctx, chartArea, scales } = chart;
						const yScale = scales.y;
						if (!yScale) return;
						isDown.forEach((down, index) => {
							if (!down) return;
							const y = yScale.getPixelForValue(index);
							ctx.save();
							ctx.fillStyle = '#DC2626';
							ctx.beginPath();
							ctx.roundRect(chartArea.left, y - 10, 6, 15, 2);
							ctx.fill();
							ctx.fillStyle = '#FFFFFF';
							ctx.font = 'bold 8px system-ui, sans-serif';
							ctx.textAlign = 'center';
							ctx.textBaseline = 'middle';
							ctx.restore();
						});
					},
				},
			],
			data: {
				labels,
				datasets: [{
					label: 'Uptime Score',
					data: scores,
					backgroundColor: (ctx) => {
						const score = scores[ctx.dataIndex ?? 0] ?? 0;
						return uptimeBarColor(score);
					},
					hoverBackgroundColor: scores.map((score) => uptimeBarColor(score)),
					borderColor: isDown.map((down) => (down ? '#991B1B' : 'transparent')),
					borderWidth: isDown.map((down) => (down ? 2 : 0)),
					borderRadius: 4,
					borderSkipped: false,
					maxBarThickness: 36,
				}],
			},
			options: {
				indexAxis: 'y',
				responsive: true,
				maintainAspectRatio: false,
				devicePixelRatio: Math.ceil(window.devicePixelRatio || 1),
				animation: { duration: 250 },
				interaction: { mode: 'index', axis: 'y', intersect: false },
				layout: { padding: { right: 12, left: 14, top: 4, bottom: 4 } },
				scales: {
					x: {
						beginAtZero: true,
						max: 100,
						border: { display: false },
						grid: { color: 'rgba(0,0,0,0.06)', drawTicks: false },
						ticks: {
							padding: 8,
							stepSize: 20,
							font: { weight: 'bold' },
							callback: (value: string | number): string => `${value}%`,
						},
					},
					y: {
						grid: { display: false },
						border: { display: false },
						ticks: {
							color: (ctx) => {
								const index = ctx.index;
								if (isDown[index]) return '#DC2626';
								const hovered = (ctx.chart as UptimeChart).$hoveredIndex;
								return ctx.index === hovered ? '#0D6EFD' : '#475569';
							},
							padding: 6,
							font: (ctx) => ({
								weight: isDown[ctx.index] ? 'bold' : 'bold',
								size: 11,
							}),
							autoSkip: labels.length > 30,
							maxTicksLimit: labels.length > 20 ? 20 : undefined,
							callback: (_value: string | number, index: number): string => {
								const label = truncateUrl(labels[index] ?? '');
								return isDown[index] ? `${label} · down` : label;
							},
						},
					},
				},
				plugins: {
					legend: { display: false },
					tooltip: {
						enabled: true,
						displayColors: true,
						padding: 10,
						titleFont: { weight: 'bold', size: 13 },
						bodyFont: { size: 12 },
						callbacks: {
							title: (items) => labels[items[0]?.dataIndex ?? 0] ?? '',
							label: (item) => {
								const score = Number(item.parsed.x);
								if (score === 0) return 'Status: Down — 0% uptime';
								return `Uptime: ${score.toFixed(1)}%`;
							},
							labelColor: (item) => {
								const score = Number(item.parsed.x);
								const color = uptimeBarColor(score);
								return {
									borderColor: color,
									backgroundColor: color,
								};
							},
						},
					},
				},
			},
		});
		const headers: additionalDataHeaders = { h1: 'Endpoint', h2: 'Uptime Score (%)' };
		return { chart, additionalData: createAdditionalData(mapped, headers, colorPicker) };
	} catch (_) {
		return { chart: emptyChart(canvas), additionalData: null };
	}
};

export const createP95Latency = (
	canvas: HTMLCanvasElement,
	data: DistributionData,
	colorPicker: ColorPicker,
): ChartData => {
	try {
		const values = data?.distribution?.values;
		if (!values || !Object.keys(values).length) throw new Error('Data is empty');
		const entries: Array<[string, number]> = Object.entries(values).sort(([, a], [, b]) => a - b);
		const mapped: Map<string, number> = new Map(entries);
		const labels: string[] = entries.map(([url]) => url);
		const latencies: number[] = entries.map(([, value]) => value);
		const colors: string[] = labels.map((url) => colorPicker.getColorForUrl(url));
		type P95Chart = Chart & { $hoveredIndex?: number | null };
		const chart = new Chart(canvas, {
			type: 'bar',
			plugins: [
				{
					id: 'yLabelHover',
					afterInit(chart) {
						const p95Chart = chart as P95Chart;
						p95Chart.$hoveredIndex = null;
						const ctx = chart.ctx;
						ctx.save();
					},
					afterEvent(chart) {
						const p95Chart = chart as P95Chart;
						const active = chart.getActiveElements();
						const nextIndex = active[0]?.index ?? null;
						if (nextIndex === p95Chart.$hoveredIndex) return;
						p95Chart.$hoveredIndex = nextIndex;
						chart.update('none');
					},
				},
			],
			data: {
				labels,
				datasets: [{
					label: 'P95 Latency',
					data: latencies,
					backgroundColor: colors,
					hoverBackgroundColor: colors,
					borderWidth: 0,
					borderRadius: 4,
					borderSkipped: false,
					maxBarThickness: 36,
				}],
			},
			options: {
				indexAxis: 'y',
				responsive: true,
				maintainAspectRatio: false,
				devicePixelRatio: Math.ceil(window.devicePixelRatio || 1),
				animation: { duration: 250 },
				interaction: { mode: 'index', axis: 'y', intersect: false },
				layout: { padding: { right: 12, left: 4, top: 4, bottom: 4 } },
				scales: {
					x: {
						beginAtZero: true,
						border: { display: false },
						grid: { color: 'rgba(0,0,0,0.06)', drawTicks: false },
						ticks: {
							padding: 8,
							font: { weight: 'bold' },
							callback: (value: string | number): string => `${value} ms`,
						},
					},
					y: {
						grid: { display: false },
						border: { display: false },
						ticks: {
							color: (ctx) => {
								const hovered = (ctx.chart as P95Chart).$hoveredIndex;
								return ctx.index === hovered ? '#0D6EFD' : '#475569';
							},
							padding: 6,
							font: { weight: 'bold', size: 11 },
							autoSkip: labels.length > 30,
							maxTicksLimit: labels.length > 20 ? 20 : undefined,
							callback: (_value: string | number, index: number): string => truncateUrl(labels[index] ?? ''),
						},
					},
				},
				plugins: {
					legend: { display: false },
					tooltip: {
						enabled: true,
						displayColors: true,
						padding: 10,
						titleFont: { weight: 'bold', size: 13 },
						bodyFont: { size: 12 },
						callbacks: {
							title: (items) => labels[items[0]?.dataIndex ?? 0] ?? '',
							label: (item) => `P95: ${Number(item.parsed.x).toLocaleString()} ms`,
						},
					},
				},
			},
		});
		const headers: additionalDataHeaders = { h1: 'Endpoint', h2: 'P95 Latency (ms)' };
		return { chart, additionalData: createAdditionalData(mapped, headers, colorPicker) };
	} catch (e) {
		return { chart: emptyChart(canvas), additionalData: null };
	}
};

export const createGeographicalTraffic = async (
	canvas: HTMLCanvasElement,
	data: any
): Promise<ChartData> => {
	try {
		const world = await fetch('https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json').then((r) => r.json());
		const countries = (topojson.feature(world, world.objects.countries) as any).features
			.filter((feature: any) => feature.properties.name !== 'Antarctica');
		const mapped = data.distribution.values ? new Map<string, number>(Object.entries(data.distribution.values)) : new Map<string, number>();
		const points: any[] = countries.map((feature: any) => {
			const value = mapped.get(feature.properties.name) ?? 0;
			return { feature, value: value > 0 ? value : null };
		});
		const all: number = points.reduce((acc: number, p: any) => acc + (p.value ?? 0), 0);
		canvas.style.backgroundColor = '#0b0f17';
		canvas.style.borderRadius = '3%';
		const chart = new ChoroplethChart(canvas, {
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
		const headers: additionalDataHeaders = { h1: 'Country', h2: 'Total Traffic' };
		return { chart, additionalData: createAdditionalData(mapped, headers) };
	} catch (_) {
		return { chart: emptyChoroplethChart(canvas), additionalData: null };
	}
}

export const createTrafficCongestionTrends = (
	canvas: HTMLCanvasElement,
	data: TrafficCongestionData,
	colorPicker: ColorPicker,
	timeframe: string,
): ChartData => {
	const getUrlCounts = (pairing: StringInt32Map | undefined): Record<string, number> =>
		pairing?.values ?? {};
	try {
		if (!data?.values?.length) throw new Error('Data is empty');
		const labels: string[] = [];
		const urlDataMap = new Map<string, number[]>();
		const cumulativeMap = new Map<string, number>();
		const totalPoints = data.values.length;
		data.values.forEach((entry: CongestionEntry, pointIndex: number) => {
			labels.push(entry.timerange);
			for (const [url, count] of Object.entries(getUrlCounts(entry.pairing))) {
				if (!urlDataMap.has(url)) {
					urlDataMap.set(url, new Array<number>(totalPoints).fill(0));
					cumulativeMap.set(url, 0);
				}
				urlDataMap.get(url)![pointIndex] = count;
				cumulativeMap.set(url, cumulativeMap.get(url)! + count);
			}
		});
		let stepSize;
		switch (timeframe) {
			case "1d":
				stepSize = 6;
				break;
			case "7d":
				stepSize = 1;
				break;
			case "30d":
				stepSize = 2;
				break;
			default:
				stepSize = 3;
		}
		const datasets = Array.from(urlDataMap.entries()).map(([url, dataArray]) => {
			const color = colorPicker.getColorForUrl(url);
			return {
				label: url,
				data: dataArray,
				backgroundColor: color,
				hoverBackgroundColor: color,
				borderWidth: 0,
				borderRadius: 0,
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
						ticks: {
							autoSkip: false,
							color: 'black',
							maxRotation: 0,
							padding: 6,
							font: { weight: 'bold' },
							callback: (_value: string | number, index: number): string =>
								index % stepSize === 0 || index === 0 || (index === labels.length - 1 && (index - 1) % stepSize !== 0) ? (labels[index] ?? '') : '',
						},
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
						filter: (item) => Number(item.parsed.y) !== 0,
						intersect: false,
						itemSort: (a, b) => (b.parsed.y as number) - (a.parsed.y as number),
						callbacks: {
							footer: function (tooltipItems) {
								const dataIndex = resolveHoveredIndex(tooltipItems, this.chart);
								if (dataIndex < 0) return '';
								const total = datasets.reduce(
									(sum, ds) => sum + (Number(ds.data[dataIndex]) || 0),
									0,
								);
								return `Total: ${total.toLocaleString()}`;
							},
							title: function (tooltipItems) {
								return formatTimeRangeTitle(labels, resolveHoveredIndex(tooltipItems, this.chart), timeframe);
							},
						},
					}
				},
			}
		});
		const headers: additionalDataHeaders = { h1: 'Endpoint', h2: 'Total' };
		return { chart, additionalData: createAdditionalData(cumulativeMap, headers, colorPicker) };
	} catch (e) {
		console.error(e)
		return { chart: emptyChart(canvas), additionalData: null };
	}
}
