<template>
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mx-4 md:mx-8 p-2 mb-16">
		<div v-for="entry in enabledMetrics" :key="entry.name">
			<Graph :name="entry.name" :data="entry.metrics" />
		</div>
	</div>
</template>

<script setup lang="ts">
import { fetchMetric } from "~/calls/dashboard";
const props = defineProps<{
	metrics: Record<string, boolean>;
}>();
const emit = defineEmits<{
	load: [value: void]
}>();

type MetricData = {
	name: string;
	metrics: any;
};
const enabledMetrics = ref<MetricData[]>([]);
const fetchAllMetrics = async (enabled: Record<string, boolean>) => {
	emit('load')
	const enabledNames = Object.keys(enabled).filter((name) => enabled[name]);
	const results = await Promise.all(enabledNames.map((name) => fetchMetric(name)));
	emit('load')
	return enabledNames.map((name, index) => ({
		name,
		metrics: results[index],
	}));
};

const loadMetrics = async () => {
	enabledMetrics.value = await fetchAllMetrics(props.metrics);
};

onMounted(loadMetrics);
watch(() => props.metrics, loadMetrics);
</script>
