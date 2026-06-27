<template>
	<div
		class="flex flex-col rounded-xl shadow-lg bg-white w-full text-black ring-1 ring-slate-100 h-96 md:h-[26rem] lg:h-[30rem]">
		<h1 class="text-2xl font-bold text-center mt-6 mb-2 shrink-0 text-slate-800">{{ props.name }}</h1>
		<div class="flex flex-col flex-1 min-h-0 px-5 pb-5 gap-3">
			<div class="relative flex-1 min-h-0 rounded-lg">
				<canvas ref="chartRef"></canvas>
			</div>
			<div v-if="hasAdditionalData" class="shrink-0 space-y-1.5">
				<p class="text-[11px] font-semibold uppercase tracking-wider text-slate-400 px-0.5">
					{{ additionalDataLabels.get(props.name) }}
				</p>
				<div ref="additionalDataRef"
					class="max-h-24 overflow-y-auto rounded-lg bg-slate-50 px-3 py-2 ring-1 ring-slate-100"></div>
			</div>
			<div class="flex justify-end shrink-0">
				<div class="relative">
					<select :value="props.timeframe"
						@change="emit('timeframeChange', { metric: props.name, timeframe: ($event.target as HTMLSelectElement).value as string })"
						class="appearance-none cursor-pointer rounded-lg border border-slate-200 bg-slate-50 py-1.5 pl-3 pr-9 text-xs font-medium text-slate-700 shadow-sm transition-colors hover:border-slate-300 focus:border-[#00F376] focus:outline-none focus:ring-2 focus:ring-[#00F376]/30">
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
</template>

<script setup lang="ts">
import { Chart, CategoryScale, LinearScale, BarController, BarElement, Tooltip } from "chart.js";
import { ChoroplethController, GeoFeature, ColorScale, ProjectionScale } from 'chartjs-chart-geo';
import { ChoroplethChart } from 'chartjs-chart-geo';
import { onMounted, toRaw } from "vue";
import { createTrafficCongestionTrends, createGeographicalTraffic, createP95Latency } from "~/composables/charts";
import { ColorPicker } from "~/composables/colorpicker";
import type { ChartData } from "~/composables/types";
const props = defineProps<{
	name: string;
	timeframe: string;
	data: any;
}>();
const emit = defineEmits<{
	timeframeChange: [{ metric: string, timeframe: string }];
}>();
const chartRef = ref<HTMLCanvasElement | null>(null);
const additionalDataRef = ref<HTMLDivElement | null>(null);
const hasAdditionalData = computed(() => ["Traffic congestion trends"].includes(props.name));
const additionalDataLabels = new Map<string, string>([["Traffic congestion trends", "Total requests by endpoint"]]);
const urls = useState<string[]>('urls');
let chartInstance: Chart | ChoroplethChart | null = null;
let colorPicker: ColorPicker | null = null;
Chart.register(CategoryScale, LinearScale, BarController, BarElement, Tooltip, ChoroplethController, GeoFeature, ColorScale, ProjectionScale);
onMounted((): void => {
	if (props.name === "Traffic congestion trends") {
		colorPicker = new ColorPicker(urls.value)
	}
	populateChart(props.data);
});
watch(() => props.data, (newData: any): void => {
	if (chartInstance) chartInstance.destroy();
	populateChart(newData);
});

const populateChart = async (data: any): Promise<void> => {
	if (props.name === "Traffic congestion trends" && chartRef.value && colorPicker) {
		const { chart, additionalData }: ChartData = createTrafficCongestionTrends(
			chartRef.value,
			toRaw(data),
			colorPicker
		);
		chartInstance = chart;
		if (additionalDataRef.value) {
			additionalDataRef.value.replaceChildren();
			if (additionalData) {
				additionalDataRef.value.appendChild(additionalData);
			}
		}
	} else if (props.name === "Geographical traffic" && chartRef.value) {
		chartInstance = await createGeographicalTraffic(chartRef.value, toRaw(data))
	} else if (props.name === "P95 Latency" && chartRef.value) {
		chartInstance = createP95Latency(chartRef.value, toRaw(data))
	}
}
</script>
