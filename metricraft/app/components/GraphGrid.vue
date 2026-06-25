<template>
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mx-4 md:mx-8 p-2 mb-16">
		<div v-for="entry in enabledMetrics" :key="entry.name">
			<Graph :name="entry.name" :data="entry.metrics" :timeframe="entry.timeframe" />
		</div>
	</div>
</template>

<script setup lang="ts">
import { fetchMetric } from "~/calls/dashboard";
const props = defineProps<{
	metrics: Record<string, { enabled: boolean, timeframe: string }>
}>();
const emit = defineEmits<{
	load: [value: void]
}>();

type MetricData = {
	name: string;
	metrics: any;
	timeframe: string;
};
const enabledMetrics = ref<MetricData[]>([]);
const fetchAllMetrics = async (enabled: Record<string, { enabled: boolean, timeframe: string }>) => {
	emit('load')
	const enabledNames = Object.keys(enabled)
	const results = await Promise.all(enabledNames.map((name) => fetchMetric(name, enabled[name].timeframe)));
	emit('load')
	return enabledNames.map((name, index) => ({
		name,
		metrics: results[index],
		timeframe: enabled[name].timeframe,
	}));
};

const loadMetrics = async () => {
	enabledMetrics.value = await fetchAllMetrics(props.metrics);
};

onMounted(loadMetrics);
watch(() => props.metrics, loadMetrics);
</script>
