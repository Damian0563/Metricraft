<template>
	<div class="flex flex-col rounded-xl shadow-lg bg-white w-full text-black ring-1 ring-slate-100"
		:class="hasAdditionalData ? 'h-96 md:h-[26rem] lg:h-[30rem]' : 'h-80 md:h-96 lg:h-108'">
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
		</div>
	</div>
</template>

<script setup lang="ts">
import { Chart, CategoryScale, LinearScale, BarController, BarElement, Tooltip } from "chart.js";
import { onMounted, toRaw } from "vue";
import { createTrafficCongestionTrends } from "~/composables/charts";
import { ColorPicker } from "~/composables/colorpicker";
import type { ChartData } from "~/composables/types";
const props = defineProps<{
	name: string;
	data: any;
}>();
const chartRef = ref<HTMLCanvasElement | null>(null);
const additionalDataRef = ref<HTMLDivElement | null>(null);
const hasAdditionalData = computed(() => props.name === "Traffic congestion trends");
const additionalDataLabels = new Map<string, string>([["Traffic congestion trends", "Total requests by endpoint"]]);
const urls = useState<string[]>('urls');
let chartInstance: Chart | null = null;
let colorPicker: ColorPicker | null = null;
Chart.register(CategoryScale, LinearScale, BarController, BarElement, Tooltip);
onMounted((): void => {
	colorPicker = new ColorPicker(urls.value)
	populateChart(props.data);
});
watch(() => props.data, (newData: any): void => {
	if (chartInstance) chartInstance.destroy();
	populateChart(newData);
});

const populateChart = (data: any): void => {
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
	}
}
</script>
