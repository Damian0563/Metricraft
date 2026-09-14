<template>
	<div class="w-full px-8 py-2 mt-4">
		<div class="flex items-center justify-between mb-3">
			<h1 class="text-3xl font-bold" style="color: #00F376;">Settings</h1>
		</div>
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<div class="lg:col-span-2 space-y-6">
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<div class="space-y-8">
						<div class="pb-6 border-b border-gray-200">
							<h2 class="text-xl font-semibold text-gray-800 mb-4">Preferences</h2>
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
						class="w-full cursor-pointer text-left px-4 py-3 rounded-xl border border-gray-200 hover:border-[#00F376] hover:shadow-md transition-all duration-300 text-gray-700 font-medium">
						Customize Dashboard View
					</button>
				</div>
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<h3 class="text-xl font-semibold text-gray-800 mb-4">Log Retention Policy</h3>
					<div class="flex flex-col gap-3">
						<div class="flex items-center justify-between px-4 py-3 rounded-xl border border-gray-100 bg-gray-50">
							<span class="text-sm font-medium text-gray-700">Current log storage</span>
							<span class="text-sm font-semibold text-gray-900">
								{{ logCapacity === null ? '—' : `${logCapacity.toFixed(2)} MB` }}
							</span>
						</div>
						<button type="button" @click="showDeleteLogs = true" :disabled="!logCapacity"
							class="w-full cursor-pointer text-left px-4 py-3 rounded-xl border border-gray-200 hover:border-red-400 hover:shadow-md transition-all duration-300 text-gray-700 font-medium disabled:opacity-50 disabled:cursor-not-allowed">
							Delete logs
						</button>
						<p v-if="errorMessage" class="text-xs text-red-500">{{ errorMessage }}</p>
					</div>
					<DeleteLogsModal :show="showDeleteLogs" :capacity="logCapacity ?? 0" @close="showDeleteLogs = false"
						@confirm="deleteLogs" />
				</div>
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<h4 class="text-xl font-semibold text-gray-800 mb-4">Configure Rules</h4>
					<button type="button"
						class="w-full cursor-pointer text-left px-4 py-3 rounded-xl border border-gray-200 hover:border-[#00F376] hover:shadow-md transition-all duration-300 text-gray-700 font-medium"
						@click="navigateTo('/rules?type=grouping')">
						Go to Rules
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { changeDerivedMetrics } from "@/calls/settings"
import { getLogCapacity, deleteLogCapacity } from "@/calls/dashboard"
type Metric = { id: number; name: string; description: string; enabled: boolean, timeframe: string }
type CompactMetric = { name: string; enabled: boolean; timeframe: string }
const props = defineProps<{
	derivedMetrics: Record<string, { enabled: boolean, timeframe: string }>;
}>();
const errorMessage = ref("")
const emit = defineEmits<{
	customizeView: [value: boolean];
	load: [value: void];
	updateMetrics: [value: CompactMetric[]];
}>();
const customizeDashboard = ref(false)
const logCapacity: Ref<number | null> = ref(null)
const showDeleteLogs = ref(false)
const pendingMetrics = ref<Metric[]>([])
const originalMetrics = ref<Metric[]>([
	{ id: 1, name: 'Geographical traffic', description: 'Map of origins of http requests in a specified time interval.', enabled: true, timeframe: "7d" },
	{ id: 2, name: 'P95 Latency', description: '95th percentile response time per endpoint', enabled: true, timeframe: "7d" },
	{ id: 3, name: 'Traffic congestion trends', description: 'Request volume measured in one hour time intervals over specified time frame.', enabled: false, timeframe: "7d" },
	{ id: 4, name: 'Uptime Score', description: 'Availability percentage over specified time frame.', enabled: true, timeframe: "7d" },
	{ id: 5, name: 'Geographic performance', description: 'Median response times broken down by clients\' countries.', enabled: false, timeframe: "7d" },
	{ id: 6, name: 'Status code distribution', description: 'Breakdown of HTTP response codes grouped by category (2xx, 3xx, 4xx, 5xx) over time.', enabled: false, timeframe: "7d" },
	{ id: 7, name: 'Route congestion', description: 'Mosty congested endpoints over time.', enabled: true, timeframe: "7d" },
	{ id: 8, name: 'Throughput', description: 'Requests per second measured over configurable time intervals to track traffic capacity and trends.', enabled: true, timeframe: "7d" },
	{ id: 9, name: 'HTTP method distribution', description: 'Breakdown of HTTP methods used by clients over time.', enabled: true, timeframe: "7d" },
	{ id: 10, name: 'Unique visitors', description: 'Number of unique visitors over time.', enabled: true, timeframe: "7d" },
	{ id: 11, name: 'Hot hours', description: 'Most popular hours of the day.', enabled: true, timeframe: "7d" },
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

const deleteLogs = async (toDeleteMbs: number) => {
	showDeleteLogs.value = false
	emit('load')
	errorMessage.value = ""
	try {
		logCapacity.value = Number(await deleteLogCapacity(toDeleteMbs))
	} catch {
		errorMessage.value = "Something went wrong, when deleting logs. Check your internet connection and try again."
	} finally {
		emit('load')
	}
}

onMounted(async () => {
	emit('load')
	try {
		logCapacity.value = Number(await getLogCapacity())
	} catch {
		errorMessage.value = "Something went wrong, when fetching log capacity. Check your internet connection and try again."
	} finally {
		emit('load')
	}
})
</script>
