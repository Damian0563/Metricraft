<template>
	<div
		class="flex flex-col rounded-xl shadow-lg bg-white w-full text-black ring-1 ring-slate-100 h-96 md:h-[26rem] lg:h-[30rem]">
		<h1 class="text-2xl font-bold text-center mt-6 mb-2 shrink-0 text-slate-800">{{ props.name }}</h1>
		<div class="flex flex-col flex-1 min-h-0 px-5 pb-5 gap-3">
			<CustomGraphHeader :data="{ metric: props.name, data: props.data }" />
			<div class="relative flex-1 min-h-0 rounded-lg">
				<canvas ref="chartRef"></canvas>
			</div>
			<div ref="additionalDataRef" class="hidden" aria-hidden="true"></div>
			<div class="flex justify-between shrink-0">
				<div class="justify-start">
					<button
						class="cursor-pointer shadow-sm rounded-lg border border-slate-200 bg-slate-50 py-1.5 pl-3 pr-9 text-xs font-medium text-slate-700 transition-colors hover:border-slate-300 focus:border-[#00F376] focus:outline-none focus:ring-2 focus:ring-[#00F376]/30"
						@click="seeDetails = !seeDetails; emit('seeDetails', { metric: props.name, additionalData: additionalDataRef })">
						View details
					</button>
				</div>
				<div class="justify-end">
					<div class="relative">
						<select :value="props.timeframe"
							@change="emit('timeframeChange', { metric: props.name, timeframe: ($event.target as HTMLSelectElement).value as string })"
							class="appearance-none cursor-pointer rounded-lg border border-slate-200 bg-slate-50 py-1.5 pl-3 pr-9 text-xs font-medium text-slate-700 shadow-sm transition-colors hover:border-slate-300 focus:border-[#00F376] focus:outline-none focus:ring-2 focus:ring-[#00F376]/30">
							<option value="1d">Last 24 hours</option>
							<option value="7d">Last 7 days</option>
							<option value="30d">Last 30 days</option>
							<option value="90d">Last 90 days</option>
							<option value="180d">Last 180 days</option>
							<option value="365d">Last 365 days</option>
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
import { Chart, CategoryScale, LinearScale, BarController, BarElement, Tooltip, PointElement, LineElement, LineController, TimeScale } from "chart.js";
import { ChoroplethController, GeoFeature, ColorScale, ProjectionScale } from 'chartjs-chart-geo';
import { ChoroplethChart } from 'chartjs-chart-geo';
import { onMounted, toRaw } from "vue";
import { useColorPicker } from "~/composables/colorpicker";
import { createTrafficCongestionTrends, createThroughput, createGeographicalTraffic, createP95Latency, createUptimeScore } from "~/composables/charts/charts";
import type { ChartData } from "~/composables/types";
const props = defineProps<{
	name: string;
	timeframe: string;
	data: any;
}>();
const emit = defineEmits<{
	timeframeChange: [{ metric: string, timeframe: string }];
	seeDetails: [{ metric: string, additionalData: HTMLDivElement | null }];
}>();
const chartRef = ref<HTMLCanvasElement | null>(null);
const additionalDataRef = ref<HTMLDivElement | null>(null);
const seeDetails = ref<boolean>(false);
const colorPicker = useColorPicker();
let chartInstance: Chart | ChoroplethChart | null = null;
Chart.register(PointElement, TimeScale, LineController, LineElement, CategoryScale, LinearScale, BarController, BarElement, Tooltip, ChoroplethController, GeoFeature, ColorScale, ProjectionScale);
onMounted((): void => {
	if (chartInstance) chartInstance.destroy();
	populateChart(props.data);
});
watch(() => props.data, (newData: any): void => {
	if (chartInstance) chartInstance.destroy();
	populateChart(newData);
});

const mutateAdditionalData = (additionalData: HTMLElement | null): void => {
	if (additionalDataRef.value) {
		additionalDataRef.value.replaceChildren();
		if (additionalData) {
			additionalDataRef.value.appendChild(additionalData);
		}
	}
}

const populateChart = async (data: any): Promise<void> => {
	const picker = colorPicker.value;
	if (props.name === "Traffic congestion trends" && chartRef.value && picker) {
		const { chart, additionalData }: ChartData = createTrafficCongestionTrends(chartRef.value, toRaw(data), picker, props.timeframe);
		chartInstance = chart;
		mutateAdditionalData(additionalData);
	} else if (props.name === "Geographical traffic" && chartRef.value) {
		const { chart, additionalData }: ChartData = await createGeographicalTraffic(chartRef.value, toRaw(data));
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
		console.log(toRaw(data))
	}
}
</script>
