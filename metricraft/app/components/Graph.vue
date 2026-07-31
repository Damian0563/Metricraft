<template>
	<div
		:class="[
			'relative flex flex-col overflow-hidden rounded-xl w-full text-black h-96 md:h-[26rem] lg:h-[30rem] transition-shadow duration-300',
			props.custom
				? 'bg-white shadow-[0_10px_40px_-8px_rgba(0,243,118,0.16),0_4px_14px_-4px_rgba(15,23,42,0.07)] ring-1 ring-[#00F376]/25'
				: 'bg-white shadow-lg ring-1 ring-slate-100',
		]">
		<div
			v-if="props.custom"
			class="pointer-events-none absolute inset-y-0 left-0 w-1 bg-gradient-to-b from-[#00F376] to-[#00B35C]"
			aria-hidden="true" />
		<div class="shrink-0 px-5 pt-6 pb-2 text-center">
			<p
				v-if="props.custom"
				class="mb-1.5 text-[10px] font-semibold uppercase tracking-[0.14em] text-[#00B35C]">
				Overwatch · Custom
			</p>
			<h1
				:class="[
					'text-2xl font-bold text-slate-800',
					props.custom && 'tracking-tight text-slate-900',
				]">
				{{ props.name }}
			</h1>
		</div>
		<div class="flex flex-col flex-1 min-h-0 px-5 pb-5 gap-3">
			<CustomGraphHeader :data="{ metric: props.name, data: props.data }" />
			<div
				:class="[
					'relative flex-1 min-h-0 rounded-lg',
					props.custom && 'bg-slate-50/70 ring-1 ring-inset ring-[#00F376]/12',
				]">
				<canvas ref="chartRef"></canvas>
			</div>
			<div ref="additionalDataRef" class="hidden" aria-hidden="true"></div>
			<div class="flex justify-between shrink-0">
				<div class="flex justify-start gap-2">
					<button
						:class="controlClass"
						@click="openDetails">
						View details
					</button>
					<button
						v-if="props.name === 'Status code distribution' || props.name === 'HTTP method distribution' || props.name === 'Traffic congestion trends'"
						:class="controlClass"
						@click="toggleDetailed">
						{{ detailedMode ? 'Grouped view' : 'Detailed view' }}
					</button>
				</div>
				<div class="justify-end">
					<div class="relative">
						<select :value="props.timeframe"
							@change="emit('timeframeChange', { metric: props.name, timeframe: ($event.target as HTMLSelectElement).value as string })"
							:class="[controlClass, 'appearance-none pr-9']">
							<option value="0.5d">Last 12 hours</option>
							<option value="1d">Last 24 hours</option>
							<option value="7d">Last 7 days</option>
							<option value="30d">Last 30 days</option>
							<option value="90d">Last 90 days</option>
							<option value="180d">Last 180 days</option>
							<option value="365d">Last 365 days</option>
							<option value="7t">This week</option>
							<option value="30t">This month</option>
							<option value="365t">This year</option>
						</select>
						<svg class="pointer-events-none absolute right-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400"
							fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
						</svg>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Chart, CategoryScale, LinearScale, BarController, BarElement, Tooltip, PointElement, LineElement, LineController, TimeScale, PieController, ArcElement, Legend } from "chart.js";
import { ChoroplethController, GeoFeature, ColorScale, ProjectionScale } from 'chartjs-chart-geo';
import { ChoroplethChart } from 'chartjs-chart-geo';
import { onMounted, toRaw } from "vue";
import { useColorPicker } from "~/composables/colorpicker";
import { createTrafficCongestionTrends, createHotHours, createRouteCongestion, createHttpMethodMix, createUniqueVisitors, createThroughput, createStatusCodeDistribution, createGeographicalTraffic, createGeographicPerformance, createP95Latency, createUptimeScore } from "~/composables/charts/charts";
import type { ChartData, WorldData } from '@/composables/types/metrics'
const props = withDefaults(defineProps<{
	name: string;
	timeframe: string;
	data: any;
	worldData: WorldData | undefined;
	custom?: boolean;
}>(), {
	custom: false,
});

const controlClass = computed(() =>
	props.custom
		? 'cursor-pointer rounded-lg border border-[#00F376]/25 bg-[#00F376]/[0.06] py-1.5 px-3 text-xs font-medium text-[#007A45] shadow-sm transition-colors hover:border-[#00F376]/40 hover:bg-[#00F376]/10 focus:border-[#00F376] focus:outline-none focus:ring-2 focus:ring-[#00F376]/30'
		: 'cursor-pointer rounded-lg border border-slate-200 bg-slate-50 py-1.5 px-3 text-xs font-medium text-slate-700 shadow-sm transition-colors hover:border-slate-300 focus:border-[#00F376] focus:outline-none focus:ring-2 focus:ring-[#00F376]/30',
);
const emit = defineEmits<{
	timeframeChange: [{ metric: string, timeframe: string }];
	seeDetails: [{ metric: string, additionalData: HTMLDivElement | null }];
}>();
const chartRef = ref<HTMLCanvasElement | null>(null);
const additionalDataRef = ref<HTMLDivElement | null>(null);
const seeDetails = ref<boolean>(false);
const colorPicker = useColorPicker();
const detailedMode = ref<boolean>(false);
let chartInstance: Chart | ChoroplethChart | null = null;
Chart.register(PointElement, TimeScale, LineController, LineElement, CategoryScale, LinearScale, BarController, BarElement, Tooltip, Legend, PieController, ArcElement, ChoroplethController, GeoFeature, ColorScale, ProjectionScale);
onMounted((): void => {
	if (chartInstance) chartInstance.destroy();
	populateChart(props.data);
});
watch(() => props.data, (newData: any): void => {
	if (chartInstance) chartInstance.destroy();
	populateChart(newData);
});

const toggleDetailed = (): void => {
	detailedMode.value = !detailedMode.value;
	if (chartInstance) chartInstance.destroy();
	populateChart(props.data);
};

const mutateAdditionalData = (additionalData: HTMLElement | null): void => {
	if (additionalDataRef.value) {
		additionalDataRef.value.replaceChildren();
		if (additionalData) {
			additionalDataRef.value.appendChild(additionalData);
		}
	}
}

const openDetails = (): void => {
	seeDetails.value = !seeDetails.value;
	emit('seeDetails', {
		metric: props.name,
		additionalData: additionalDataRef.value,
	});
}

const populateChart = async (data: any): Promise<void> => {
	const picker = colorPicker.value;
	if (props.name === "Traffic congestion trends" && chartRef.value && picker) {
		const { chart, additionalData }: ChartData = createTrafficCongestionTrends(chartRef.value, toRaw(data), picker, props.timeframe, detailedMode.value);
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "Geographical traffic" && chartRef.value) {
		const { chart, additionalData }: ChartData = createGeographicalTraffic(chartRef.value, toRaw(data), props.worldData);
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "P95 Latency" && chartRef.value && picker) {
		const { chart, additionalData }: ChartData = createP95Latency(chartRef.value, toRaw(data), picker)
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "Uptime Score" && chartRef.value && picker) {
		const { chart, additionalData }: ChartData = createUptimeScore(chartRef.value, toRaw(data), picker)
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "Throughput" && chartRef.value) {
		const { chart, additionalData }: ChartData = createThroughput(chartRef.value, toRaw(data), props.timeframe)
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "Geographic performance" && chartRef.value) {
		const { chart, additionalData }: ChartData = createGeographicPerformance(chartRef.value, toRaw(data), props.worldData)
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "Status code distribution" && chartRef.value) {
		const { chart, additionalData }: ChartData = createStatusCodeDistribution(chartRef.value, toRaw(data), detailedMode.value)
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "Route congestion" && chartRef.value && picker) {
		const { chart, additionalData }: ChartData = createRouteCongestion(chartRef.value, toRaw(data), picker)
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "HTTP method distribution" && chartRef.value) {
		const { chart, additionalData }: ChartData = createHttpMethodMix(chartRef.value, toRaw(data), props.timeframe, detailedMode.value)
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "Unique visitors" && chartRef.value) {
		const { chart, additionalData }: ChartData = createUniqueVisitors(chartRef.value, toRaw(data), props.timeframe)
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "Hot hours" && chartRef.value) {
		const { chart, additionalData }: ChartData = createHotHours(chartRef.value, toRaw(data))
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	}
}
</script>
