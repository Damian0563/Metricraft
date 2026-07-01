<template>
	<div class="w-full px-8 py-2 ">
		<div class="flex items-center justify-between mb-3">
			<h1 class="text-3xl font-bold" style="color: #00F376;">Settings</h1>
		</div>
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<div class="lg:col-span-2 space-y-6">
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<div class="space-y-8">
						<div class="pb-6 border-b border-gray-200">
							<h2 class="text-xl font-semibold text-gray-800 mb-4">Preferences</h2>
							<label class="flex items-center justify-between cursor-pointer">
								<div class="flex items-center gap-2">
									<span class="text-base font-medium text-gray-700">Real-time updates</span>
									<div class="group relative">
										<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#00F376"
											class="w-5 h-5 text-gray-400">
											<path fill-rule="evenodd"
												d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12Zm8.706-1.442c1.146-.573 2.437.463 2.126 1.706l-.709 2.836.042-.02a.75.75 0 01.67 1.34l-.04.022c-1.147.573-2.438-.463-2.127-1.706l.71-2.836-.042.02a.75.75 0 11-.671-1.34l.041-.022ZM12 9a.75.75 0 100-1.5.75.75 0 000 1.5Z"
												clip-rule="evenodd" />
										</svg>
										<div
											class="absolute left-0 bottom-full mb-1 hidden group-hover:block w-72 p-3 bg-[#00F376] text-black text-s rounded-lg shadow-lg z-10">
											<div> This setting is highly discouraged if your service expects heavy traffic.
											</div>
										</div>
									</div>
								</div>
								<div class="relative">
									<input type="checkbox" class="sr-only peer" :checked="realtimeEnabled"
										@change="emit('realtimeToggle', ($event.target as HTMLInputElement).checked)" />
									<div
										class="w-11 h-6 bg-gray-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#00F376]">
									</div>
								</div>
							</label>
						</div>
						<div>
							<h2 class="text-xl font-semibold text-gray-800 mb-4">Derived Metrics</h2>
							<div class="space-y-3">
								<div v-for="metric in pendingMetrics" :key="metric.id"
									class="flex items-center justify-between px-4 py-3 rounded-xl border border-gray-100 bg-gray-50">
									<div>
										<p class="text-sm font-medium text-gray-700">{{ metric.name }}</p>
										<p class="text-xs text-gray-500">{{ metric.description }}</p>
									</div>
									<div class="flex items-center gap-3">
										<label class="relative cursor-pointer">
											<input type="checkbox" class="sr-only peer" v-model="metric.enabled" />
											<div
												class="w-11 h-6 bg-gray-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#00F376]">
											</div>
										</label>
										<span :class="[
											'text-xs font-medium px-3 py-1 rounded-full',
											metric.enabled ? 'bg-[#00F376]/20 text-green-700' : 'bg-gray-200 text-gray-500'
										]">
											{{ metric.enabled ? 'Active' : 'Inactive' }}
										</span>
									</div>
								</div>
							</div>
							<button @click="applyMetricChanges"
								class="mt-4 px-6 py-2 bg-[#00F376] text-gray-900 font-semibold rounded-lg hover:bg-[#00D96A] transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
								:disabled="!hasChanges">
								Apply Changes
							</button>
						</div>
					</div>
				</div>
			</div>
			<div class="space-y-6">
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<h2 class="text-xl font-semibold text-gray-800 mb-4">Team</h2>
					<span @click="navigateTo(`/invite`)"
						class="text-base font-medium text-gray-700 hover:cursor-pointer hover:text-[#00F376] transition-colors duration-200">
						Manage team members
					</span>
				</div>
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<h2 class="text-xl font-semibold text-gray-800 mb-4">Customization</h2>
					<button @click="emit('customizeView', !customizeDashboard)"
						class="w-full text-left px-4 py-3 rounded-xl border border-gray-200 hover:border-[#00F376] hover:shadow-md transition-all duration-300 text-gray-700 font-medium">
						Customize Dashboard View
					</button>
				</div>
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<h2 class="text-xl font-semibold text-gray-800 mb-4">Log Retention Policy</h2>
					<div class="flex flex-col gap-3">
						<select v-model="logRetention"
							@change="changeRetention(Number(logRetention)); emit('changeRetention', Number(logRetention))"
							class="px-4 py-2 rounded-lg border border-gray-200 text-gray-700 font-medium focus:outline-none focus:border-[#00F376] transition-colors duration-200">
							<option value="7">7 days</option>
							<option value="30">30 days</option>
							<option value="90">90 days</option>
							<option value="180">6 months</option>
							<option value="365">1 year</option>
						</select>
						<p class="text-sm text-gray-500">Automatically delete logs older than the selected period to reduce
							memory usage. The data of the derived metrics will be compacted and still available, but raw http
							traffic logs will be deleted, you can export them at any time.</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { changeDerivedMetrics, changeRetention } from "@/calls/settings"
type Metric = { id: number; name: string; description: string; enabled: boolean, timeframe: string }
type CompactMetric = { name: string; enabled: boolean; timeframe: string }
const props = defineProps<{
	realtimeEnabled: boolean;
	logRetention: number;
	derivedMetrics: Record<string, { enabled: boolean, timeframe: string }>;
}>();
const emit = defineEmits<{
	realtimeToggle: [value: boolean];
	customizeView: [value: boolean];
	load: [value: void];
	updateMetrics: [value: CompactMetric[]];
	changeRetention: [value: Number];
}>();
const customizeDashboard = ref(false)
const logRetention = ref(props.logRetention)
watch(() => props.logRetention, (val) => logRetention.value = val)
const pendingMetrics = ref<Metric[]>([])
const originalMetrics = ref<Metric[]>([
	{ id: 1, name: 'Geographical traffic', description: 'Map of origins of http requests in a specified time interval.', enabled: true, timeframe: "7d" },
	{ id: 2, name: 'P95 Latency', description: '95th percentile response time per endpoint', enabled: true, timeframe: "7d" },
	{ id: 3, name: 'Traffic congestion trends', description: 'Request volume measured in one hour time intervals over specified time frame.', enabled: false, timeframe: "7d" },
	{ id: 4, name: 'Uptime Score', description: 'Availability percentage over specified time frame.', enabled: true, timeframe: "7d" },
	{ id: 5, name: 'Geographic performance', description: 'Average response times and error rates broken down by client country or region.', enabled: false, timeframe: "7d" },
	{ id: 6, name: 'Status code distribution', description: 'Breakdown of HTTP response codes grouped by category (2xx, 3xx, 4xx, 5xx) over time.', enabled: false, timeframe: "7d" },
	{ id: 7, name: 'Median response time', description: 'P50 latency across all requests, providing a representative measure of typical endpoint performance.', enabled: true, timeframe: "7d" },
	{ id: 8, name: 'Throughput', description: 'Requests per second measured over configurable time intervals to track traffic capacity and trends.', enabled: true, timeframe: "7d" },
])
watch(() => props.derivedMetrics, (metrics) => {
	const updated = originalMetrics.value.map(metric => {
		const enabled = metrics[metric.name]
		return typeof enabled === "object" ? { ...metric, enabled: enabled.enabled, timeframe: enabled.timeframe } : { ...metric, enabled: false, timeframe: "7d" }
	})
	originalMetrics.value = updated
	pendingMetrics.value = updated.map(m => ({ ...m }))
}, { immediate: true, deep: true })

const hasChanges = computed(() =>
	originalMetrics.value.some((orig, i) => orig.enabled !== pendingMetrics.value[i]?.enabled)
)
const applyMetricChanges = async () => {
	emit('load')
	if (!hasChanges.value) return
	const changes = pendingMetrics.value.map(m => ({ name: m.name, enabled: m.enabled, timeframe: m.timeframe }))
	await changeDerivedMetrics(changes)
	originalMetrics.value = pendingMetrics.value.map(m => ({ ...m }))
	emit('updateMetrics', changes)
	emit('load')
}
</script>
