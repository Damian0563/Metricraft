import { Chart } from "chart.js";
import { emptyChart } from "@/composables/charts/charts";
import { createAdditionalData } from "@/composables/charts/chartUtils";
import { truncateUrl } from "@/composables/helpers";
import { ColorPicker } from "~/composables/colorpicker";
import type { ChartData, additionalDataHeaders } from '@/composables/types/metrics'
import type { ChartType } from '@/composables/types/additional'

const formatGranularTooltipTitle = (
	labels: string[],
	dataIndex: number,
	timeframe: string,
): string => {
	if (dataIndex < 0 || dataIndex >= labels.length) return '';
	if (timeframe === '1d' || timeframe === '7d') {
		const start = labels[dataIndex];
		const end = labels[dataIndex + 1];
		if (timeframe === '1d') return end != null ? `${start} - ${end}` : `${start}-${labels[0]}`;
		return end != null ? `${start} - ${end}` : (start ?? '');
	}
	return labels[dataIndex] ?? '';
};

const stepSizeForTimeframe = (timeframe: string, labelCount: number): number => {
	switch (timeframe) {
		case '1d':
			return 6;
		case '7d':
			return 1;
		case '30d':
			return 2;
		default:
			return labelCount <= 8 ? 1 : labelCount <= 16 ? 2 : 3;
	}
};

const resolveHoveredIndex = (
	tooltipItems: Array<{ dataIndex?: number }>,
	chartInstance: Chart,
): number =>
	tooltipItems.find((item) => item.dataIndex != null)?.dataIndex ??
	chartInstance.getActiveElements()[0]?.index ??
	-1;

const xTickLabel = (
	labels: string[],
	index: number,
	stepSize: number,
): string =>
	index % stepSize === 0 || index === 0 || (index === labels.length - 1 && (index - 1) % stepSize !== 0)
		? (labels[index] ?? '')
		: '';

export const genericAccumulatedChart = (
	canvas: HTMLCanvasElement,
	data: Array<{ grouping: string, value: number }> | undefined,
): ChartData => {
	try {
		if (!data || data.length === 0) return { chart: emptyChart(canvas), additionalData: null };
		const entries = data.filter((entry) => entry.grouping.trim() !== "");
		if (!entries.length) return { chart: emptyChart(canvas), additionalData: null };
		const sorted = [...entries].sort((a, b) => b.value - a.value);
		const sortedMap = new Map(sorted.map(({ grouping, value }) => [grouping, value]));
		const labels: string[] = sorted.map(({ grouping }) => grouping);
		const values: number[] = sorted.map(({ value }) => value);
		const colorPicker = new ColorPicker(labels);
		const colors: string[] = labels.map((label) => colorPicker?.getColorForInstance(label) ?? '#94A3B8');
		const total: number = values.reduce((sum, value) => sum + value, 0);
		if (total <= 0) return { chart: emptyChart(canvas), additionalData: null };
		const chart: Chart = new Chart(canvas, {
			type: 'pie',
			data: {
				labels,
				datasets: [{
					label: 'Total',
					data: values,
					backgroundColor: colors,
					hoverBackgroundColor: colors,
					borderColor: '#FFFFFF',
					borderWidth: 2,
					hoverOffset: 8,
				}],
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				devicePixelRatio: Math.ceil(window.devicePixelRatio || 1),
				animation: { duration: 250 },
				layout: { padding: { top: 8, right: 8, bottom: 8, left: 8 } },
				plugins: {
					legend: {
						display: true,
						position: 'right',
						labels: {
							usePointStyle: true,
							pointStyle: 'circle',
							boxWidth: 8,
							boxHeight: 8,
							padding: 12,
							color: '#475569',
							font: { weight: 'bold', size: 11 },
							generateLabels: (chartInstance) => {
								const dataset = chartInstance.data.datasets[0];
								return (chartInstance.data.labels ?? []).map((label, index) => {
									const color = Array.isArray(dataset?.backgroundColor)
										? dataset.backgroundColor[index]
										: dataset?.backgroundColor;
									return {
										text: truncateUrl(String(label)),
										fillStyle: color as string,
										strokeStyle: color as string,
										lineWidth: 0,
										hidden: !chartInstance.getDataVisibility(index),
										index,
									};
								});
							},
						},
					},
					tooltip: {
						enabled: true,
						displayColors: true,
						padding: 10,
						titleFont: { weight: 'bold', size: 13 },
						bodyFont: { size: 12 },
						callbacks: {
							title: (items) => labels[items[0]?.dataIndex ?? 0] ?? '',
							label: (item) => {
								const value = Number(item.parsed);
								const share = total > 0 ? ((value / total) * 100).toFixed(1) : '0.0';
								return `Total: ${value.toLocaleString()} (${share}%)`;
							},
							labelColor: (item) => {
								const color = colors[item.dataIndex ?? 0] ?? '#475569';
								return { borderColor: color, backgroundColor: color };
							},
						},
					},
				},
			},
		});
		const headers: additionalDataHeaders = { h1: 'Grouping', h2: 'Total' };
		return { chart, additionalData: createAdditionalData(sortedMap, headers, colorPicker) };
	} catch (_) {
		return { chart: emptyChart(canvas), additionalData: null };
	}
}

export const genericGranularChart = (
	canvas: HTMLCanvasElement,
	data: Array<{ grouping: string, value: number }> | undefined,
	chartType: ChartType | undefined,
	timeframe: string = '',
): ChartData => {
	const barFill = '#00F376';
	const barHover = '#00C962';
	const barBorder = '#009E4F';
	const lineFill = '#00F376';
	const lineStroke = '#00C962';
	try {
		if (!data || data.length === 0) return { chart: emptyChart(canvas), additionalData: null };
		const entries = data.filter((entry) => entry.grouping.trim() !== "");
		if (!entries.length) return { chart: emptyChart(canvas), additionalData: null };
		const labels: string[] = entries.map(({ grouping }) => grouping);
		const values: number[] = entries.map(({ value }) => value);
		const mapped = entries.map(({ grouping, value }) => ({ timerange: grouping, value }));
		const stepSize = stepSizeForTimeframe(timeframe, labels.length);
		const isBar = chartType === 'bar';
		type GranularChart = Chart & { $hoveredIndex?: number | null };
		const chart: Chart = new Chart(canvas, {
			type: isBar ? 'bar' : 'line',
			plugins: [
				{
					id: 'genericGranularLabelHover',
					afterInit(chart) {
						const granularChart = chart as GranularChart;
						granularChart.$hoveredIndex = null;
					},
					afterEvent(chart) {
						const granularChart = chart as GranularChart;
						const active = chart.getActiveElements();
						const nextIndex = active[0]?.index ?? null;
						if (nextIndex === granularChart.$hoveredIndex) return;
						granularChart.$hoveredIndex = nextIndex;
						chart.update('none');
					},
				},
			],
			data: {
				labels,
				datasets: [{
					label: 'Value',
					data: values,
					...(isBar ? {
						backgroundColor: barFill,
						hoverBackgroundColor: barHover,
						borderColor: barBorder,
						borderWidth: 1,
						borderRadius: 6,
						borderSkipped: false,
						maxBarThickness: 40,
						categoryPercentage: 0.82,
						barPercentage: 0.9,
					} : {
						showLine: true,
						tension: 0.3,
						borderColor: lineStroke,
						borderWidth: 2,
						backgroundColor: lineFill,
						pointBackgroundColor: lineFill,
						pointBorderColor: '#FFFFFF',
						pointBorderWidth: 2,
						pointRadius: 4,
						pointHoverRadius: 6,
						pointHoverBackgroundColor: lineFill,
						pointHoverBorderColor: '#FFFFFF',
						pointHoverBorderWidth: 2,
					}),
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
							color: (ctx) => {
								const hovered = (ctx.chart as GranularChart).$hoveredIndex;
								return ctx.index === hovered ? '#0D6EFD' : '#475569';
							},
							font: (ctx): any => ({
								weight: ctx.index === (ctx.chart as GranularChart).$hoveredIndex ? 'bold' : '600',
								size: 11,
							}),
							callback: (_value: string | number, index: number): string =>
								xTickLabel(labels, index, stepSize),
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
						ticks: {
							padding: 8,
							precision: 0,
							color: '#475569',
							font: { weight: 'bold', size: 11 },
						},
						title: {
							display: true,
							text: 'Value',
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
							title(tooltipItems) {
								return formatGranularTooltipTitle(
									labels,
									resolveHoveredIndex(tooltipItems, this.chart),
									timeframe,
								);
							},
							label: (item) => `Value: ${Number(item.parsed.y).toLocaleString()}`,
							labelColor: () => ({
								borderColor: isBar ? barFill : lineFill,
								backgroundColor: isBar ? barFill : lineFill,
							}),
						},
					},
				},
			},
		});
		const addData = timeframe === '1d'
			? values.map((value, index) => ({
				timerange: formatGranularTooltipTitle(labels, index, timeframe),
				value,
			}))
			: mapped;
		const headers: additionalDataHeaders = { h1: 'Timerange', h2: 'Value' };
		return { chart, additionalData: createAdditionalData(addData, headers) };
	} catch (e) {
		console.error(e)
		return { chart: emptyChart(canvas), additionalData: null };
	}
}
