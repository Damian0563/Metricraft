<template>
	<div>
		<ClientOnly>
			<Spinner :loading="loading" />
		</ClientOnly>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<UpdateMetricForm :open="editingMetric !== null" :metric="editingMetric" :saving="loading"
			@close="editingMetric = null" @update="onMetricUpdate" />
		<div class="relative flex items-center justify-center mb-6">
			<h1 class="text-3xl font-bold text-center" style="color: #00F376;">Overwatch</h1>
		</div>
		<div class="max-w-8xl mx-auto grid gap-4 lg:grid-cols-[minmax(0,1fr)_40rem] lg:items-stretch">
			<div class="flex flex-col gap-4 min-w-0">
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<h2 class="text-xl font-semibold text-gray-800 mb-4">What is Overwatch?</h2>
					<p class="text-sm text-gray-600 leading-relaxed mb-4">
						Overwatch lets you define <span class="font-semibold text-gray-800">custom metrics</span> from
						the shape of your own API. Instead of only tracking generic HTTP stats, you describe how your
						endpoints look — request body, headers, query parameters — and Overwatch extracts exactly the
						values you care about.
					</p>
					<p class="text-sm text-gray-600 leading-relaxed">
						Point a metric at a JSON field, a header, or a query param and choose how it should be
						aggregated. Metricraft then charts it on your dashboards alongside everything else.
					</p>
				</div>
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<div class="flex items-start gap-3 mb-6">
						<div
							class="h-10 w-10 shrink-0 rounded-full bg-[#00F376]/10 flex items-center justify-center text-[#00B35C]">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"
								aria-hidden="true">
								<path fill-rule="evenodd"
									d="M12.316 3.051a1 1 0 01.633 1.265l-4 12a1 1 0 11-1.898-.632l4-12a1 1 0 011.265-.633zM5.707 6.293a1 1 0 010 1.414L3.414 10l2.293 2.293a1 1 0 11-1.414 1.414l-3-3a1 1 0 010-1.414l3-3a1 1 0 011.414 0zm8.586 0a1 1 0 011.414 0l3 3a1 1 0 010 1.414l-3 3a1 1 0 11-1.414-1.414L16.586 10l-2.293-2.293a1 1 0 010-1.414z"
									clip-rule="evenodd" />
							</svg>
						</div>
						<div class="min-w-0">
							<h2 class="text-xl font-semibold text-gray-800">Create a custom metric</h2>
							<p class="text-sm text-gray-500 mt-1">
								Describe where the value lives in your API traffic and how to measure it.
							</p>
						</div>
					</div>
					<div class="space-y-6">
						<div>
							<div class="flex items-baseline justify-between gap-3 mb-2">
								<label for="metric-name" class="text-sm font-medium text-gray-700">Metric name</label>
								<p v-if="namingError" class="text-xs text-red-500 text-right">{{ namingError }}</p>
							</div>
							<input id="metric-name" v-model="metricName" type="text" placeholder="Checkout order total"
								class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors" />
							<p class="mt-2 text-xs text-gray-500">A human-friendly label shown on your dashboards.</p>
						</div>
						<div class="flex flex-col gap-8 xl:flex-row xl:justify-between">
							<div class="shrink-0">
								<label for="metric-method" class="block text-sm font-medium text-gray-700 mb-2">Method</label>
								<select id="metric-method" v-model="method"
									class="w-full max-w-xs px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors cursor-pointer">
									<option v-for="verb in methods" :key="verb" :value="verb">{{ verb }}</option>
								</select>
							</div>
							<div class="flex-1 min-w-0">
								<AutoSuggestedUrlInput v-model="path" input-id="metric-path" label="Endpoint path"
									help-text="Overwatch inspects requests matching this route." @pattern-error="pathError = $event"
									@submit-event="addMetric" />
							</div>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">Extract value from</label>
							<div class="grid grid-cols-3 gap-2">
								<button v-for="src in sources" :key="src.id" type="button" @click="source = src.id" :class="[
									'flex flex-col items-center gap-1.5 px-3 py-3 rounded-lg border text-xs font-semibold transition-colors cursor-pointer',
									source === src.id
										? 'border-[#00F376] bg-[#00F376]/10 text-[#00B35C] ring-1 ring-[#00F376]/30'
										: 'border-gray-200 text-gray-600 hover:border-gray-300 hover:bg-gray-50']">
									<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"
										v-html="src.icon" />
									{{ src.label }}
								</button>
							</div>
						</div>
						<div>
							<label for="metric-selector" class="block text-sm font-medium text-gray-700 mb-2">
								{{ selectorLabel }}
							</label>
							<input id="metric-selector" v-model="selector" type="text" :placeholder="selectorPlaceholder"
								class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 font-mono text-sm focus:outline-none focus:border-[#00F376] transition-colors" />
							<p class="mt-2 text-xs text-gray-500">{{ selectorHelp }}</p>
						</div>

						<div class="flex flex-col gap-8 xl:flex-row xl:justify-between">
							<div class="flex-1 min-w-0">
								<label for="metric-type" class="block text-sm font-medium text-gray-700 mb-2">Value type</label>
								<select id="metric-type" v-model="valueType"
									class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors cursor-pointer">
									<option v-for="t in valueTypes" :key="t" :value="t">{{ t }}</option>
								</select>
							</div>
							<div class="flex-1 min-w-0">
								<label for="metric-agg" class="block text-sm font-medium text-gray-700 mb-2">Aggregation</label>
								<select id="metric-agg" v-model="aggregation"
									class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors cursor-pointer">
									<option v-for="agg in allowedAggregations" :key="agg" :value="agg">{{ agg }}</option>
								</select>
							</div>
							<div class="flex-1 min-w-0">
								<label for="metric-timeframe" class="block text-sm font-medium text-gray-700 mb-2">Timeframe</label>
								<select id="metric-timeframe" v-model="timeframe"
									class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors cursor-pointer">
									<option v-for="[label, value] in timeframes" :key="value" :value="value">{{ label }}</option>
								</select>
							</div>
							<div class="flex-1 min-w-0 ">
								<label for="metric-chart-type" class="block text-sm font-medium text-gray-700 mb-2">Chart type</label>
								<select id="metric-chart-type" v-model="chartType"
									class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors cursor-pointer">
									<option v-for="ct in chartTypes" :key="ct" :value="ct">{{ ct }}</option>
								</select>
							</div>
						</div>
						<div class="pt-6 border-t border-gray-100 flex justify-between gap-3">
							<div>
								<button type="button" @click="addMetric" :disabled="!canSave"
									class="px-8 py-3 bg-[#00F376] text-gray-900 font-bold rounded-lg hover:bg-[#00D96A] transition-all shadow-lg disabled:opacity-50 disabled:cursor-not-allowed uppercase tracking-wider text-sm cursor-pointer">
									Create metric
								</button>
								<label for="apply-rules" class="inline-flex items-center gap-2 cursor-pointer shrink-0 mx-5">
									<span class="select-none text-sm font-medium text-gray-700">Apply grouping rules</span>
									<div class="group relative mb-4">
										<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#00F376"
											class="w-5 h-5 text-gray-400">
											<path fill-rule="evenodd"
												d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12Zm8.706-1.442c1.146-.573 2.437.463 2.126 1.706l-.709 2.836.042-.02a.75.75 0 01.67 1.34l-.04.022c-1.147.573-2.438-.463-2.127-1.706l.71-2.836-.042.02a.75.75 0 11-.671-1.34l.041-.022ZM12 9a.75.75 0 100-1.5.75.75 0 000 1.5Z"
												clip-rule="evenodd" />
										</svg>
										<div v-if="source !== 'query'"
											class="absolute right-0 bottom-full mb-2 hidden group-hover:block w-96 p-3 bg-[#00F376] text-black text-s rounded-lg shadow-lg z-10">
											Toggle this on if the above {{ path }} is a grouping metric. For example if the path is a subpath
											and you want to inspect traffic on all of the subpaths of the request i.e <span
												class="font-semibold text-gray-800">{{ path }}/id</span> then you should enable this option.
											custom metrics. This does not apply to query params.
										</div>
									</div>
									<div class="relative">
										<input type="checkbox" id="apply-rules" v-model="applyRules" class="sr-only peer"
											:disabled="source === 'query'" />
										<div :style="{ cursor: source === 'query' ? 'not-allowed' : 'default' }"
											class="w-9 h-5 bg-gray-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#00F376]">
										</div>
									</div>
								</label>
							</div>
							<button type="button" @click="resetForm"
								class="px-6 py-3 bg-gray-100 text-gray-700 font-semibold rounded-lg hover:bg-gray-200 transition-colors cursor-pointer">
								Reset
							</button>
						</div>
					</div>
				</div>
			</div>

			<div
				class="bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden lg:sticky lg:top-8 flex flex-col">
				<section class="border-b border-gray-100">
					<div class="px-6 py-5 border-b border-gray-100">
						<h2 class="text-xl font-semibold text-gray-800">Request preview</h2>
						<p class="text-sm text-gray-500 mt-1">
							A simulated request for this configuration. The
							<span class="font-semibold text-[#00B35C]">highlighted</span> part is what Overwatch inspects.
						</p>
					</div>
					<div class="p-4">
						<div class="rounded-lg bg-gray-900 text-gray-100 text-xs leading-relaxed p-4 overflow-x-auto font-mono">
							<div class="whitespace-pre"><span class="text-[#00F376] font-semibold">{{ method }}</span><span>{{ ' ' +
								requestBasePath }}</span><span v-if="querySuffix"
									class="rounded bg-[#00F376]/20 ring-1 ring-[#00F376]/50 px-0.5">{{ querySuffix }}</span><span
									class="text-gray-500"> HTTP/1.1</span></div>
							<div v-for="header in headerLines" :key="header.name" class="whitespace-pre rounded"
								:class="header.target ? 'bg-[#00F376]/20 ring-1 ring-[#00F376]/50 px-0.5' : ''">
								<span class="text-sky-300">{{ header.name }}</span><span class="text-gray-400">: </span><span
									class="text-emerald-300">{{ header.value }}</span><span v-if="header.target"
									class="text-[#00F376]/70"> ← inspected</span>
							</div>
							<template v-if="showBody">
								<div class="whitespace-pre">&nbsp;</div>
								<div v-for="(line, i) in bodyLines" :key="`b-${i}`" class="whitespace-pre rounded"
									:class="line.target ? 'bg-[#00F376]/20 ring-1 ring-[#00F376]/50 px-0.5 text-emerald-300' : ''">{{
										line.text }}<span v-if="line.target" class="text-[#00F376]/70"> ← inspected</span></div>
							</template>
						</div>
						<div class="mt-3 rounded-lg border border-gray-200 bg-gray-50 px-4 py-3">
							<p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1.5">Overwatch will
								compute</p>
							<div class="flex flex-wrap items-center gap-2 text-xs">
								<span
									class="inline-flex items-center px-2.5 py-1 rounded-full bg-[#00F376]/10 text-[#00B35C] font-bold border border-[#00F376] uppercase tracking-wide">
									{{ aggregation }}
								</span>
								<span class="text-gray-500">of</span>
								<span
									class="inline-flex items-center gap-1 px-2 py-1 rounded-md bg-white border border-gray-200 text-gray-700 font-medium">
									{{ sourceLabel }}
								</span>
								<span v-if="selector" class="font-mono text-gray-800 font-semibold">'s {{ selector }}</span>
								<span v-else class="text-gray-400 italic">no field selected yet</span>
								<span class="text-gray-500">as a</span>
								<span class="font-medium text-gray-700">{{ valueType }}</span>
								<span class="text-gray-500">over</span>
								<span class="font-medium text-gray-700">{{ timeframeLabelFor(timeframe).toLowerCase() }}</span>
								<span class="text-gray-500">as a</span>
								<span class="font-medium text-gray-700">{{ chartType }}</span>
								<span class="text-gray-500">chart</span>
								<template v-if="applyRules">
									<span class="text-gray-500">with</span>
									<span class="font-medium text-gray-700">grouping rules applied</span>
								</template>
							</div>
						</div>
					</div>
				</section>
				<section class="flex flex-col flex-1 min-h-0 bg-gray-50/60">
					<div class="px-6 py-5 border-b border-gray-100">
						<div class="flex items-center justify-between gap-3">
							<h2 class="text-sm font-semibold text-gray-800">Configured metrics</h2>
							<span class="shrink-0 text-xs font-medium text-gray-500">{{ metrics.length }} defined</span>
						</div>
					</div>

					<div v-if="metrics.length > 0" class="min-h-0 flex-1 overflow-y-auto">
						<ClientOnly>
							<AnimatePresence mode="popLayout">
								<motion.div v-for="(metric, index) in metrics" :key="metricKey(metric)" layout
									class="group relative overflow-hidden border-b border-gray-100 last:border-b-0"
									:initial="{ opacity: 0, height: 0 }" :animate="{ opacity: 1, height: 'auto' }"
									:exit="{ opacity: 0, height: 0, x: 32, transition: { duration: 0.22, ease: 'easeIn' } }"
									:transition="{ type: 'spring', stiffness: 480, damping: 34, delay: Math.min(index * 0.045, 0.27) }">
									<div
										class="flex items-center justify-between gap-3 px-6 py-4 transition-colors duration-200 hover:bg-[#00F376]/[0.04]">
										<motion.div
											class="min-w-0 flex-1 border-l-2 border-transparent pl-3 transition-[border-color] duration-200 group-hover:border-[#00F376]/70"
											:initial="{ opacity: 0, x: -10 }" :animate="{ opacity: 1, x: 0 }"
											:transition="{ delay: Math.min(index * 0.045, 0.27) + 0.08, duration: 0.28, ease: 'easeOut' }">
											<p class="text-sm font-semibold text-gray-900 transition-colors group-hover:text-[#00B35C]">
												{{ metric.name }}
											</p>
											<p class="mt-1 font-mono text-xs text-gray-600 break-all">
												<span class="text-[#00B35C] font-semibold">{{ metric.method }}</span>
												{{ ' ' + metric.path }}
											</p>
											<p class="mt-2 text-xs text-gray-500">
												<span class="font-medium text-gray-700">{{ metric.aggregation }}</span>
												of {{ metric.source }} <span class="font-mono text-gray-700">{{ metric.selector }}</span>
												· {{ metric.valueType }} · {{ timeframeLabelFor(metric.timeframe).toLowerCase() }} · {{
													metric.chartType
												}}
											</p>
											<p v-if="metric.lastUpdate" class="mt-1 text-xs text-gray-400">
												Last updated {{ formatMetricLastUpdate(metric.lastUpdate) }}
											</p>
										</motion.div>
										<motion.div class="flex shrink-0 items-center gap-3" :initial="{ opacity: 0, scale: 0.9 }"
											:animate="{ opacity: 1, scale: 1 }"
											:transition="{ type: 'spring', stiffness: 420, damping: 26, delay: Math.min(index * 0.045, 0.27) + 0.12 }">
											<button type="button" @click="editingMetric = metric"
												class="shrink-0 px-4 py-2 text-sm font-semibold rounded-lg bg-[#00F376] text-gray-900 shadow-sm hover:bg-[#00D96A] hover:shadow-md transition-all cursor-pointer">
												Update
											</button>
											<button type="button"
												class="shrink-0 px-4 py-2 text-sm font-semibold rounded-lg bg-red-50 text-red-600 hover:bg-red-100 transition-colors cursor-pointer"
												@click="deleteMetric(metric)">Delete</button>
										</motion.div>
									</div>
								</motion.div>
							</AnimatePresence>
						</ClientOnly>
					</div>

					<div v-else class="min-h-0 flex-1 px-6 py-10 text-center">
						<div class="mx-auto mb-3 h-10 w-10 rounded-full bg-gray-100 flex items-center justify-center text-gray-400">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd"
									d="M3 3a1 1 0 000 2h14a1 1 0 100-2H3zm0 4a1 1 0 000 2h9a1 1 0 100-2H3zm0 4a1 1 0 100 2h14a1 1 0 100-2H3zm0 4a1 1 0 100 2h9a1 1 0 100-2H3z"
									clip-rule="evenodd" />
							</svg>
						</div>
						<p class="text-sm text-gray-500">No custom metrics yet.</p>
						<p class="text-xs text-gray-400 mt-1">Create a metric to start tracking API values.</p>
					</div>
				</section>
			</div>
		</div>
	</div>
</template>


<script setup lang="ts">
import type { ChartType, CustomMetric, MetricSource } from '@/composables/types/additional'
import { customMetricNamingError, timeframeLabelFor } from '@/composables/helpers'
import { addCustomMetric, getCustomMetrics, deleteCustomMetric, updateCustomMetric } from '@/calls/overwatch'
import { motion, AnimatePresence } from 'motion-v'
definePageMeta({
	layout: 'dashboard',
})
const metrics = ref<CustomMetric[]>([])
const loading = ref(false)
const errorMessage = ref('')
const editingMetric = ref<CustomMetric | null>(null)
const metricKey = (metric: CustomMetric) => `${metric.name}-${metric.path}-${metric.selector}-${metric.aggregation}-${metric.timeframe}-${metric.valueType}-${metric.chartType}-${metric.applyRules}`
const formatMetricLastUpdate = (iso: string | Date) => {
	const parsed = typeof iso === 'string' ? new Date(iso) : iso
	if (Number.isNaN(parsed.getTime())) return ''
	const tz = Intl.DateTimeFormat().resolvedOptions().timeZone
	const parts = new Intl.DateTimeFormat('en-US', {
		timeZone: tz,
		month: 'short',
		day: 'numeric',
		year: 'numeric',
		hour: 'numeric',
		minute: '2-digit',
		hour12: true,
	}).formatToParts(parsed)
	const get = (type: Intl.DateTimeFormatPartTypes) => parts.find(p => p.type === type)?.value ?? ''
	return `${get('month')} ${get('day')}, ${get('year')} at ${get('hour')}:${get('minute')}${get('dayPeriod').toLowerCase()}`
}

const { data: fetchedMetrics, error: errorFetchMetrics } = await useAsyncData('getCustomMetrics', () => getCustomMetrics(), { default: () => [] as CustomMetric[] })
watch(fetchedMetrics, (val) => {
	metrics.value = val ? [...val] : []
}, { immediate: true })
if (errorFetchMetrics.value) {
	errorMessage.value = 'Failed to fetch existing custom metrics.'
}
type BodyLine = { text: string; target: boolean }
const metricName = ref('')
const applyRules = ref(false)
const chartType = ref<ChartType>('line')
const chartTypes = ['line', 'bar', 'pie'] as const
const method = ref<string>('POST')
const path = ref('')
const pathError = ref('')
const source = ref<MetricSource>('body')
const selector = ref('')
const aggregation = ref<string>('avg')
const sourceLabel = computed(() => sources.find(s => s.id === source.value)?.label ?? '')
const methods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'] as const
const valueType: Ref<string> = ref('number')
const aggregationTypesPerValueType: { [key: string]: string[] } = {
	'number': ['count', 'sum', 'avg', 'min', 'max', 'p95', 'unique'],
	'string': ['count', 'unique'],
	'boolean': ['count', 'p50'],
}
const aggregationTypesPerChartType: { [key: string]: string[] } = {
	'line': ['count', 'sum', 'avg', 'min', 'max', 'p95', 'unique'],
	'bar': ['count', 'sum', 'avg', 'min', 'max', 'p95', 'unique'],
	'pie': ['count', 'unique'],
}
const valueTypes = ['number', 'string', 'boolean'] as const
const allowedAggregations = computed(() => {
	const byValue = aggregationTypesPerValueType[valueType.value] ?? []
	const byChart = aggregationTypesPerChartType[chartType.value] ?? []
	return byValue.filter(x => byChart.includes(x))
})
watch([valueType, chartType], () => {
	const allowed = allowedAggregations.value
	if (!allowed.includes(aggregation.value)) {
		aggregation.value = allowed[0] ?? 'count'
	}
}, { immediate: true })

const namingError = computed(() => customMetricNamingError(metricName.value))
const timeframes = new Map<string, string>([['Last 12 hours', "0.5d"], ['Last 24 hours', "1d"], ['Last 7 days', "7d"], ['Last 30 days', "30d"], ['Last 90 days', "90d"], ['Last 180 days', "180d"], ['Last 365 days', "365d"], ['This week', '7t'], ['This month', "30t"], ['This year', "365t"]])
const defaultTimeframe = '7d'
const timeframe = ref(defaultTimeframe)
const sources: { id: MetricSource; label: string; icon: string }[] = [
	{
		id: 'body',
		label: 'Body',
		icon: '<path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 011 1v3a1 1 0 01-1 1H4a1 1 0 01-1-1V4zm0 6a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zm10 0a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" clip-rule="evenodd" />',
	},
	{
		id: 'header',
		label: 'Header',
		icon: '<path fill-rule="evenodd" d="M3 3a1 1 0 000 2h14a1 1 0 100-2H3zm0 4a1 1 0 000 2h9a1 1 0 100-2H3zm0 4a1 1 0 100 2h14a1 1 0 100-2H3zm0 4a1 1 0 100 2h9a1 1 0 100-2H3z" clip-rule="evenodd" />',
	},
	{
		id: 'query',
		label: 'Query param',
		icon: '<path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />',
	},
]


const selectorLabel = computed(() => {
	switch (source.value) {
		case 'body': return 'JSON field path'
		case 'header': return 'Header name'
		case 'query': return 'Query parameter name'
	}
})
const selectorPlaceholder = computed(() => {
	switch (source.value) {
		case 'body': return 'order.total'
		case 'header': return 'X-Response-Time'
		case 'query': return 'page'
	}
})
const selectorHelp = computed(() => {
	switch (source.value) {
		case 'body': return 'Dot-notation path into the JSON request body, e.g. items[0].price.'
		case 'header': return 'The HTTP header whose value Overwatch should read.'
		case 'query': return 'The query string parameter Overwatch should read.'
	}
})
const sampleLeaf = computed<{ cls: string; display: string }>(() => {
	switch (valueType.value) {
		case 'number': return { cls: 'text-emerald-300', display: '42.50' }
		case 'boolean': return { cls: 'text-emerald-300', display: 'true' }
		default: return { cls: 'text-amber-200', display: '"sample"' }
	}
})

const hasBody = computed(() => ['POST', 'PUT', 'PATCH'].includes(method.value))
const showBody = computed(() => source.value === 'body' || hasBody.value)
const requestBasePath = computed(() => path.value.trim() || '/your/endpoint')
const querySuffix = computed(() =>
	source.value === 'query' && selector.value.trim()
		? `?${selector.value.trim()}=<value>`
		: '')
interface HeaderLine { name: string; value: string; target: boolean }
const headerLines = computed<HeaderLine[]>(() => {
	const lines: HeaderLine[] = [
		{ name: 'Host', value: 'api.yourservice.com', target: false },
	]
	if (showBody.value) {
		lines.push({ name: 'Content-Type', value: 'application/json', target: false })
	}
	if (source.value === 'header' && selector.value.trim()) {
		const name = selector.value.trim()
		const existing = lines.find(l => l.name.toLowerCase() === name.toLowerCase())
		if (existing) {
			existing.target = true
		} else {
			lines.push({ name, value: '<value>', target: true })
		}
	}
	return lines
})
const TARGET_TOKEN = '__OVERWATCH_TARGET__'
const buildBodyObject = (segments: string[]): unknown => {
	const [head, ...rest] = segments
	if (head === undefined) return TARGET_TOKEN
	if (/^\d+$/.test(head)) return [buildBodyObject(rest)]
	return { [head]: buildBodyObject(rest) }
}

const bodyLines = computed<BodyLine[]>(() => {
	const selecting = source.value === 'body'
	const bodyObject = selecting && selector.value.trim() !== ''
		? buildBodyObject(selector.value.trim().split('.'))
		: { orderId: 'a1b2c3', total: 42.5, currency: 'USD' }

	return JSON.stringify(bodyObject, null, 2).split('\n').map((line) => {
		const target = line.includes(TARGET_TOKEN)
		return {
			text: target ? line.replace(`"${TARGET_TOKEN}"`, sampleLeaf.value.display) : line,
			target,
		}
	})
})

const canSave = computed(() =>
	metricName.value.trim() !== '' && path.value.trim() !== '' && pathError.value === '' && namingError.value === '' && selector.value.trim() !== '')

const addMetric = async () => {
	if (!canSave.value) return
	loading.value = true
	try {
		const metric: Omit<CustomMetric, 'lastUpdate'> = {
			name: metricName.value.trim(),
			method: method.value,
			path: path.value.trim(),
			source: source.value,
			selector: selector.value.trim(),
			timeframe: timeframe.value,
			aggregation: aggregation.value,
			valueType: valueType.value,
			applyRules: applyRules.value,
			chartType: chartType.value,
		}
		await addCustomMetric(metric)
		errorMessage.value = 'Custom metric saved successfully.'
		metrics.value.push({ ...metric, lastUpdate: new Date().toISOString() })
		resetForm()
	} catch (e) {
		console.error(e)
		errorMessage.value = 'Failed to save custom metric.'
	} finally {
		loading.value = false
	}
}

const onMetricUpdate = (updated: CustomMetric) => {
	const original = editingMetric.value
	if (!original) return
	updateMetric(original, updated)
}

const updateMetric = async (original: CustomMetric, updated: CustomMetric) => {
	editingMetric.value = null
	if (metricKey(original) === metricKey(updated)) return
	loading.value = true
	try {
		await updateCustomMetric(original, updated)
		errorMessage.value = 'Custom metric updated successfully.'
		const index = metrics.value.findIndex(m => metricKey(m) === metricKey(original))
		if (index !== -1) {
			metrics.value[index] = { ...updated, lastUpdate: new Date().toISOString() }
		}
	} catch (e) {
		console.error(e)
		errorMessage.value = 'Failed to update custom metric.'
	} finally {
		loading.value = false
	}
}

const deleteMetric = async (metric: CustomMetric) => {
	try {
		await deleteCustomMetric(metric)
		errorMessage.value = 'Custom metric deleted successfully.'
		metrics.value.splice(metrics.value.findIndex(m => m === metric), 1)
	} catch (e) {
		console.error(e)
		errorMessage.value = 'Failed to delete custom metric.'
	}
}

const resetForm = () => {
	metricName.value = ''
	method.value = 'POST'
	path.value = ''
	source.value = 'body'
	selector.value = ''
	aggregation.value = 'avg'
	valueType.value = 'number'
	applyRules.value = false
	chartType.value = 'line'
	timeframe.value = defaultTimeframe
}
</script>
