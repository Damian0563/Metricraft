import { Chart } from "chart.js";
import { emptyChart } from "@/composables/charts/charts";
import { createAdditionalData } from "@/composables/charts/chartUtils";
import type { ChartData, additionalDataHeaders } from '@/composables/types/metrics';
import type { GenericChartData } from '@/composables/types/additional';
import { ColorPicker } from "~/composables/colorpicker";

const TICK_COLOR = '#475569';
const HOVER_ACCENT = '#0D6EFD';
const LINE_BORDER = '#00C962';
const LINE_POINT = '#00F376';
const BAR_FILL = '#00F376';
const BAR_HOVER = '#00C962';
const BAR_BORDER = '#009E4F';
const FALLBACK_COLORS = ['#00F376', '#38BDF8', '#A78BFA', '#FBBF24', '#F87171', '#FB923C', '#34D399', '#94A3B8'];

type HoveredChart = Chart & { $hoveredIndex?: number | null };

const labelHoverPlugin = (id: string) => ({
	id,
	afterInit(chart: Chart) {
		(chart as HoveredChart).$hoveredIndex = null;
	},
	afterEvent(chart: Chart) {
		const hoveredChart = chart as HoveredChart;
		const nextIndex = chart.getActiveElements()[0]?.index ?? null;
		if (nextIndex === hoveredChart.$hoveredIndex) return;
		hoveredChart.$hoveredIndex = nextIndex;
		chart.update('none');
	},
});

const parseData = (data: GenericChartData) => {
	if (!data?.length) throw new Error('Data is empty');
	const labels: string[] = [];
	const values: number[] = [];
	for (const entry of data) {
		if (entry.grouping === '' || (!entry.value && entry.value !== 0)) continue;
		labels.push(entry.grouping);
		values.push(Number(entry.value));
	}
	if (!labels.length) throw new Error('Data is empty');
	const mapped = new Map(labels.map((label, index) => [label, values[index]!]));
	return { labels, values, mapped };
};

const sliceColor = (label: string, index: number, colorPicker?: ColorPicker): string =>
	colorPicker?.getColorForInstance(label) ?? FALLBACK_COLORS[index % FALLBACK_COLORS.length]!;

const tickStepSize = (timeframe: string, labelCount: number): number => {
	switch (timeframe) {
		case '1d': return 6;
		case '7d': return 1;
		case '30d': return 2;
		default:
			return labelCount <= 8 ? 1 : labelCount <= 16 ? 2 : 3;
	}
};

const sharedChartOptions = {
	responsive: true,
	maintainAspectRatio: false,
	devicePixelRatio: Math.ceil(window.devicePixelRatio || 1),
	animation: { duration: 250 },
};

const sharedTooltip = {
	enabled: true,
	displayColors: true,
	padding: 10,
	titleFont: { weight: 'bold' as const, size: 13 },
	bodyFont: { size: 12 },
};

const axisTitle = (text: string) => ({
	display: true,
	text,
	color: TICK_COLOR,
	font: { weight: 'bold' as const, size: 13 },
});

export const genericAccumulatedChart = (
	canvas: HTMLCanvasElement,
	data: GenericChartData,
): ChartData => {
	try {
		const { labels, values, mapped } = parseData(data);
		const colorPicker = new ColorPicker(labels);
		const colors = labels.map((label, index) => sliceColor(label, index, colorPicker));
		const total = values.reduce((sum, value) => sum + value, 0);
		if (total <= 0) throw new Error('Data is empty');

		const chart = new Chart(canvas, {
			type: 'pie',
			data: {
				labels,
				datasets: [{
					label: 'Value',
					data: values,
					backgroundColor: colors,
					hoverBackgroundColor: colors,
					borderColor: '#FFFFFF',
					borderWidth: 2,
					hoverOffset: 8,
				}],
			},
			options: {
				...sharedChartOptions,
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
							color: TICK_COLOR,
							font: { weight: 'bold', size: 11 },
						},
					},
					tooltip: {
						...sharedTooltip,
						callbacks: {
							title: (items) => labels[items[0]?.dataIndex ?? 0] ?? '',
							label: (item) => {
								const value = Number(item.parsed);
								const share = total > 0 ? ((value / total) * 100).toFixed(1) : '0.0';
								return `Value: ${value.toLocaleString()} (${share}%)`;
							},
							labelColor: (item) => {
								const color = colors[item.dataIndex ?? 0] ?? FALLBACK_COLORS[0]!;
								return { borderColor: color, backgroundColor: color };
							},
						},
					},
				},
			},
		});

		const headers: additionalDataHeaders = { h1: 'Group', h2: 'Value' };
		return { chart, additionalData: createAdditionalData(mapped, headers, colorPicker ?? null) };
	} catch (e) {
		console.error(e);
		return { chart: emptyChart(canvas), additionalData: null };
	}
};

export const genericGranularChart = (
	canvas: HTMLCanvasElement,
	data: GenericChartData,
	chartType: 'line' | 'bar' | undefined,
	timeframe = '',
): ChartData => {
	try {
		const { labels, values, mapped } = parseData(data);
		const stepSize = tickStepSize(timeframe, labels.length);
		const pluginId = chartType === 'bar' ? 'genericBarLabelHover' : 'genericLineLabelHover';
		const isBar = chartType === 'bar';

		const chart = new Chart(canvas, {
			type: isBar ? 'bar' : 'line',
			plugins: [labelHoverPlugin(pluginId)],
			data: {
				labels,
				datasets: [{
					label: 'Value',
					data: values,
					...(isBar ? {
						backgroundColor: BAR_FILL,
						hoverBackgroundColor: BAR_HOVER,
						borderColor: BAR_BORDER,
						borderWidth: 1,
						borderRadius: 6,
						borderSkipped: false,
						maxBarThickness: 40,
						categoryPercentage: 0.82,
						barPercentage: 0.9,
					} : {
						showLine: true,
						tension: 0.3,
						borderColor: LINE_BORDER,
						borderWidth: 2,
						backgroundColor: LINE_POINT,
						pointBackgroundColor: LINE_POINT,
						pointBorderColor: '#FFFFFF',
						pointBorderWidth: 2,
						pointRadius: 4,
						pointHoverRadius: 6,
						pointHoverBackgroundColor: LINE_POINT,
						pointHoverBorderColor: '#FFFFFF',
						pointHoverBorderWidth: 2,
					}),
				}],
			},
			options: {
				...sharedChartOptions,
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
								const hovered = (ctx.chart as HoveredChart).$hoveredIndex;
								return ctx.index === hovered ? HOVER_ACCENT : TICK_COLOR;
							},
							font: (ctx): any => ({
								weight: ctx.index === (ctx.chart as HoveredChart).$hoveredIndex ? 'bold' : '600',
								size: 11,
							}),
							callback: (_value: string | number, index: number): string =>
								index % stepSize === 0 || index === 0 || (index === labels.length - 1 && (index - 1) % stepSize !== 0)
									? (labels[index] ?? '')
									: '',
						},
						title: { ...axisTitle('Time'), padding: { top: 2 } },
					},
					y: {
						beginAtZero: true,
						border: { display: false },
						grid: { color: 'rgba(0,0,0,0.06)', drawTicks: false },
						ticks: {
							padding: 8,
							precision: 0,
							color: TICK_COLOR,
							font: { weight: 'bold', size: 11 },
						},
						title: axisTitle('Value'),
					},
				},
				plugins: {
					legend: { display: false },
					tooltip: {
						...sharedTooltip,
						callbacks: {
							title: (items) => labels[items[0]?.dataIndex ?? 0] ?? '',
							label: (item) => `Value: ${Number(item.parsed.y).toLocaleString()}`,
							labelColor: () => ({
								borderColor: isBar ? BAR_FILL : LINE_POINT,
								backgroundColor: isBar ? BAR_FILL : LINE_POINT,
							}),
						},
					},
				},
			},
		});
		const headers: additionalDataHeaders = { h1: 'Timerange', h2: 'Value' };
		return { chart, additionalData: createAdditionalData(mapped, headers) };
	} catch (e) {
		console.error(e);
		return { chart: emptyChart(canvas), additionalData: null };
	}
};
