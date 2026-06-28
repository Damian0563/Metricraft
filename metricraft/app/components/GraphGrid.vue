<template>
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mx-4 md:mx-8 p-2 mb-16">
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<AdditionalData :data="additionalData" v-if="viewingDetails" @close="viewingDetails = false" />
		<div v-for="entry in enabledMetrics" :key="entry.name">
			<Graph :name="entry.name" :data="entry.metrics" :timeframe="entry.timeframe"
				@timeframe-change="handleTimeframeChange($event)" @see-details="handleDetails($event)" />
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
const enabledMetrics = ref<MetricData[] | undefined>([]);
const errorMessage = ref<string>('');
const viewingDetails = ref<boolean>(false);
const additionalData = ref<HTMLDivElement | null>(null);
const fetchAllMetrics = async (enabled: Record<string, { enabled: boolean, timeframe: string }>) => {
	emit('load')
	try {
		const enabledEntries = Object.entries(enabled).filter(
			(entry): entry is [string, { enabled: boolean, timeframe: string }] => entry[1]?.enabled === true
		)
		const results = await Promise.all(enabledEntries.map(([name, config]) => fetchMetric(name, config.timeframe)))
		emit('load')
		return enabledEntries.map(([name, config], index) => ({
			name,
			metrics: results[index],
			timeframe: config.timeframe,
		}))
	} catch (_) {
		emit('load')
		errorMessage.value = 'Something went wrong, Check your internet connection and try again.'
	}
};

const handleDetails = (div: HTMLDivElement | null) => {
	if (div) {
		additionalData.value = div;
		viewingDetails.value = true;
	}
}

const handleTimeframeChange = async (obj: { metric: string, timeframe: string }) => {
	const { metric, timeframe } = obj;
	if (!enabledMetrics.value) return;
	emit('load');
	try {
		const metrics = await fetchMetric(metric, timeframe, true);
		enabledMetrics.value = enabledMetrics.value.map((entry) =>
			entry.name === metric ? { ...entry, metrics, timeframe } : entry
		);
	} catch (_) {
		errorMessage.value = 'Something went wrong, Check your internet connection and try again.';
	} finally {
		emit('load');
	}
};

const loadMetrics = async () => {
	enabledMetrics.value = await fetchAllMetrics(props.metrics);
};

onMounted(loadMetrics);
watch(() => props.metrics, loadMetrics);
</script>
