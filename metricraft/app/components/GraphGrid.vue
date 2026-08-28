<template>
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mx-4 md:mx-8 p-2 mb-16">
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<AdditionalData :show="viewingDetails" :data="additionalData" :metric="additionalDataName"
			@close="viewingDetails = false" />
		<ClientOnly>
			<AnimatePresence>
				<motion.div v-for="(entry, i) in enabledMetrics" :key="entry.name" layout :initial="{ opacity: 0, height: 0 }"
					:animate="{ opacity: 1, height: 'auto' }"
					:exit="{ opacity: 0, height: 0, x: 32, transition: { duration: 0.22, ease: 'easeIn' } }"
					:transition="{ type: 'spring', stiffness: 480, damping: 34, delay: Math.min(i * 0.045, 0.27) }">
					<Graph :name="entry.name" :custom="!!entry.customMetrics" :accumulate="!!entry.accumulate"
						:definition="entry.definition ?? null" :data="entry.metrics" :timeframe="entry.timeframe"
						:worldData="worldData" @timeframe-change="handleTimeframeChange($event)"
						@see-details="handleDetails($event)" @metric-updated="loadMetrics"
						@error="errorMessage = $event" />
				</motion.div>
			</AnimatePresence>
		</ClientOnly>
	</div>
</template>

<script setup lang="ts">
import { fetchMetric, fetchCustomMetrics } from "~/calls/dashboard";
import { topojson } from 'chartjs-chart-geo';
import { timeframeValueFor } from '@/composables/helpers';
import type { MetricData } from '@/composables/types/additional';
import type { WorldData } from '@/composables/types/metrics';
import { motion, AnimatePresence } from 'motion-v';
const props = defineProps<{
	metrics: Record<string, { enabled: boolean, timeframe: string }>
}>();
const emit = defineEmits<{
	load: [value: void]
}>();
const enabledMetrics = ref<MetricData[] | undefined>([]);
const errorMessage = ref<string>('');
const viewingDetails = ref<boolean>(false);
const additionalData = shallowRef<HTMLDivElement | null>(null);
const additionalDataName = ref<string>('');
const { data: worldData } = await useFetch('https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json', {
	transform: (world: any): WorldData => ({
		countries: (topojson.feature(world, world.objects.countries) as any).features
			.filter((feature: any) => feature.properties.name !== 'Antarctica'),
	}),
});
const fetchAllMetrics = async (enabled: Record<string, { enabled: boolean, timeframe: string }>) => {
	emit('load')
	try {
		const enabledEntries = Object.entries(enabled).filter(
			(entry): entry is [string, { enabled: boolean, timeframe: string }] => entry[1]?.enabled === true
		)
		const results = await Promise.all(enabledEntries.map(([name, config]) => fetchMetric(name, config.timeframe, errorMessage)))
		const customMetrics: MetricData[] = await fetchCustomMetrics(errorMessage)
		const standardMetrics: MetricData[] = enabledEntries.map(([name, config], index) => ({
			name,
			metrics: results[index],
			timeframe: config.timeframe,
			customMetrics: false,
		}));
		return standardMetrics.concat(customMetrics.map((metric) => ({ ...metric, timeframe: timeframeValueFor(metric.timeframe) })));
	} catch (_) {
		errorMessage.value = 'Something went wrong, Check your internet connection and try again.'
	} finally {
		emit('load')
	}
};

const handleDetails = (event: { metric: string, additionalData: HTMLDivElement | null }) => {
	const { metric, additionalData: detailsEl } = event;
	if (detailsEl) {
		additionalData.value = detailsEl;
		additionalDataName.value = metric;
		viewingDetails.value = true;
	}
}

const handleTimeframeChange = async (obj: { metric: string, timeframe: string }) => {
	const { metric, timeframe } = obj;
	if (!enabledMetrics.value) return;
	emit('load');
	try {
		const metrics = await fetchMetric(metric, timeframe, errorMessage, true);
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
