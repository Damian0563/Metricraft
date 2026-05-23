<template>
	<div class="flex flex-col rounded shadow-lg bg-white w-full h-80 md:h-96 lg:h-108 text-black">
		<h1 class="text-2xl font-bold text-center mt-6 mb-2 shrink-0">{{ props.name }}</h1>
		<div class="relative flex-1 min-h-0 p-6">
			<canvas ref="chartRef"></canvas>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Chart, CategoryScale, LinearScale, BarController, BarElement, Tooltip } from "chart.js";
import { onMounted, toRaw } from "vue";
import { createTrafficCongestionTrends } from "~/composables/charts";
import { ColorPicker } from "~/composables/colorpicker";
const props = defineProps<{
	name: string;
	data: any;
}>();
const chartRef = ref<HTMLCanvasElement | null>(null);
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
		chartInstance = createTrafficCongestionTrends(chartRef.value, toRaw(data), colorPicker);
	}
}
</script>
