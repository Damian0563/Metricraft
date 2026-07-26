<template>
	<div class="flex flex-col w-full text-black h-56">
		<h3 class="text-sm font-semibold text-slate-800 break-all shrink-0 mb-2">{{ props.url }}</h3>
		<div class="relative flex-1 min-h-0 rounded-lg">
			<p v-if="displayEmptyMessage" class="text-m text-gray-500 text-center text-decoration-none">No data available yet.
			</p>
			<canvas ref="chartRef"></canvas>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Chart, CategoryScale, LinearScale, BarController, BarElement, Tooltip, Title } from "chart.js";
import { onBeforeUnmount, onMounted, toRaw } from "vue";
import { createWorkerUptimeChart } from "~/composables/charts/uptime";
import type { WorkerUptimeData } from '@/composables/types/additional'
const props = defineProps<{
	url: string;
	data: WorkerUptimeData;
}>();
const chartRef = ref<HTMLCanvasElement | null>(null);
const displayEmptyMessage = ref(false);
let chartInstance: Chart | null = null;
Chart.register(CategoryScale, LinearScale, BarController, BarElement, Tooltip, Title);
const renderChart = (data: WorkerUptimeData): void => {
	if (data.entries?.length === 0) {
		displayEmptyMessage.value = true;
		return;
	}
	displayEmptyMessage.value = false;
	if (chartInstance) chartInstance.destroy();
	if (!chartRef.value) return;
	chartInstance = createWorkerUptimeChart(chartRef.value, toRaw(data));
};
onMounted(() => renderChart(props.data));
watch(() => props.data, (newData: WorkerUptimeData) => renderChart(newData));
onBeforeUnmount(() => {
	if (chartInstance) chartInstance.destroy();
});
</script>
