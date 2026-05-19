<template>
	<div class="flex flex-col rounded shadow-lg bg-white w-full h-48 md:h-64 lg:h-96 text-black aspect-square">
		<h1 class="text-2xl font-bold text-center mt-2">{{ props.name }}</h1>
		<canvas ref="chartRef" class="p-12"></canvas>
	</div>
</template>

<script setup lang="ts">
import { Chart, CategoryScale, LinearScale, BarController, BarElement, Tooltip } from "chart.js";
import { onMounted, toRaw } from "vue";
import { createTrafficCongestionTrends } from "~/composables/charts";
const props = defineProps<{
	name: string;
	data: any;
}>();
const chartRef = ref<HTMLCanvasElement | null>(null);
let chartInstance: Chart | null = null;
Chart.register(CategoryScale, LinearScale, BarController, BarElement, Tooltip);
onMounted((): void => {
	populateChart(props.data);
});
watch(() => props.data, (newData: any): void => {
	if (chartInstance) chartInstance.destroy();
	populateChart(newData);
});

const populateChart = (data: any): void => {
	if (props.name === "Traffic congestion trends" && chartRef.value) {
		chartInstance = createTrafficCongestionTrends(chartRef.value, toRaw(data));
	}
}
</script>
