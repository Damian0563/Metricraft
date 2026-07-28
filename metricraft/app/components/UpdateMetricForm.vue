<template>
	<ClientOnly>
		<Teleport to="body">
			<AnimatePresence>
				<motion.div v-if="open" key="update-metric-overlay"
					class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm px-4"
					:initial="{ opacity: 0 }" :animate="{ opacity: 1 }" :exit="{ opacity: 0 }" :transition="{ duration: 0.25 }"
					@click.self="emit('close')">
					<motion.div
						class="relative flex w-full max-w-3xl max-h-[95vh] flex-col overflow-hidden rounded-xl bg-white shadow-xl"
						:initial="{ opacity: 0, scale: 0.9, y: 24 }" :animate="{ opacity: 1, scale: 1, y: 0 }"
						:exit="{ opacity: 0, scale: 0.9, y: 24 }" :transition="{ type: 'spring', duration: 0.45, bounce: 0.25 }"
						@click.stop>
						<div class="shrink-0 border-b border-gray-100 px-8 py-6">
							<div class="flex items-start justify-between gap-3">
								<div>
									<h2 class="text-xl font-semibold text-gray-800">Update metric</h2>
									<p class="mt-1 text-sm text-gray-500">Change how Overwatch extracts and aggregates this value.</p>
								</div>
								<button type="button" @click="emit('close')"
									class="shrink-0 rounded-lg p-1.5 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 cursor-pointer"
									aria-label="Close">
									<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
										<path fill-rule="evenodd"
											d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
											clip-rule="evenodd" />
									</svg>
								</button>
							</div>
						</div>

						<div class="min-h-0 flex-1 overflow-y-auto px-8 py-6">
							<div class="space-y-6">
								<div>
									<label :for="`${idPrefix}-name`" class="mb-1.5 block text-sm font-medium text-gray-700">Metric
										name</label>
									<input :id="`${idPrefix}-name`" v-model="metricName" type="text" placeholder="Checkout order total"
										class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-800 transition-colors focus:border-[#00F376] focus:outline-none" />
								</div>

								<div class="flex flex-col gap-4 sm:flex-row">
									<div class="shrink-0 sm:w-36">
										<label :for="`${idPrefix}-method`"
											class="mb-1.5 block text-sm font-medium text-gray-700">Method</label>
										<select :id="`${idPrefix}-method`" v-model="method"
											class="w-full cursor-pointer rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-800 transition-colors focus:border-[#00F376] focus:outline-none">
											<option v-for="verb in methods" :key="verb" :value="verb">{{ verb }}</option>
										</select>
									</div>
									<div class="min-w-0 flex-1">
										<AutoSuggestedUrlInput v-model="path" :input-id="`${idPrefix}-path`" label="Endpoint path"
											help-text="Overwatch inspects requests matching this route." @pattern-error="pathError = $event"
											@submit-event="submit" />
									</div>
								</div>

								<div>
									<label class="mb-1.5 block text-sm font-medium text-gray-700">Extract value from</label>
									<div class="grid grid-cols-3 gap-2">
										<button v-for="src in sources" :key="src.id" type="button" @click="source = src.id" :class="[
											'flex cursor-pointer flex-col items-center gap-1 rounded-lg border px-2 py-2 text-xs font-semibold transition-colors',
											source === src.id
												? 'border-[#00F376] bg-[#00F376]/10 text-[#00B35C] ring-1 ring-[#00F376]/30'
												: 'border-gray-200 text-gray-600 hover:border-gray-300 hover:bg-gray-50',
										]">
											<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"
												v-html="src.icon" />
											{{ src.label }}
										</button>
									</div>
								</div>

								<div>
									<label :for="`${idPrefix}-selector`" class="mb-1.5 block text-sm font-medium text-gray-700">
										{{ selectorLabel }}
									</label>
									<input :id="`${idPrefix}-selector`" v-model="selector" type="text" :placeholder="selectorPlaceholder"
										class="w-full rounded-lg border border-gray-200 px-3 py-2 font-mono text-sm text-gray-800 transition-colors focus:border-[#00F376] focus:outline-none" />
									<p class="mt-1.5 text-xs text-gray-500">{{ selectorHelp }}</p>
								</div>

								<div class="grid gap-4 sm:grid-cols-2">
									<div>
										<label :for="`${idPrefix}-type`" class="mb-1.5 block text-sm font-medium text-gray-700">Value
											type</label>
										<select :id="`${idPrefix}-type`" v-model="valueType"
											class="w-full cursor-pointer rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-800 transition-colors focus:border-[#00F376] focus:outline-none">
											<option v-for="t in valueTypes" :key="t" :value="t">{{ t }}</option>
										</select>
									</div>
									<div>
										<label :for="`${idPrefix}-agg`"
											class="mb-1.5 block text-sm font-medium text-gray-700">Aggregation</label>
										<select :id="`${idPrefix}-agg`" v-model="aggregation"
											class="w-full cursor-pointer rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-800 transition-colors focus:border-[#00F376] focus:outline-none">
											<option v-for="agg in aggregationTypes[valueType]" :key="agg" :value="agg">{{ agg }}</option>
										</select>
									</div>
									<div>
										<label :for="`${idPrefix}-timeframe`"
											class="mb-1.5 block text-sm font-medium text-gray-700">Timeframe</label>
										<select :id="`${idPrefix}-timeframe`" v-model="timeframe"
											class="w-full cursor-pointer rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-800 transition-colors focus:border-[#00F376] focus:outline-none">
											<option v-for="tf in timeframes" :key="tf" :value="tf">{{ tf }}</option>
										</select>
									</div>
									<div>
										<label :for="`${idPrefix}-chart-type`" class="mb-1.5 block text-sm font-medium text-gray-700">Chart
											type</label>
										<select :id="`${idPrefix}-chart-type`" v-model="chartType"
											class="w-full cursor-pointer rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-800 transition-colors focus:border-[#00F376] focus:outline-none">
											<option v-for="ct in chartTypes" :key="ct" :value="ct">{{ ct }}</option>
										</select>
									</div>
								</div>

								<label :for="`${idPrefix}-apply-rules`" class="inline-flex cursor-pointer items-center gap-2">
									<span class="select-none text-sm font-medium text-gray-700">Apply grouping rules</span>
									<div class="relative">
										<input :id="`${idPrefix}-apply-rules`" v-model="applyRules" type="checkbox" class="peer sr-only"
											:disabled="source === 'query'" />
										<div :style="{ cursor: source === 'query' ? 'not-allowed' : 'default' }"
											class="h-5 w-9 rounded-full bg-gray-300 after:absolute after:start-[2px] after:top-[2px] after:h-4 after:w-4 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-[#00F376] peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none">
										</div>
									</div>
								</label>
							</div>
						</div>

						<div class="shrink-0 flex justify-end gap-3 border-t border-gray-100 px-8 py-5">
							<button type="button" @click="emit('close')" :disabled="saving"
								class="cursor-pointer rounded-lg bg-gray-100 px-6 py-3 text-sm font-semibold text-gray-700 transition-colors hover:bg-gray-200 disabled:cursor-not-allowed disabled:opacity-50">
								Cancel
							</button>
							<button type="button" @click="submit" :disabled="!canSave || saving"
								class="cursor-pointer rounded-lg bg-[#00F376] px-8 py-3 text-sm font-bold uppercase tracking-wider text-gray-900 shadow-lg transition-all hover:bg-[#00D96A] disabled:cursor-not-allowed disabled:opacity-50">
								Save changes
							</button>
						</div>
					</motion.div>
				</motion.div>
			</AnimatePresence>
		</Teleport>
	</ClientOnly>
</template>

<script setup lang="ts">
import { AnimatePresence, motion } from 'motion-v'
import type { ChartType, CustomMetric, MetricSource } from '@/composables/types/additional'

const props = defineProps<{
	open: boolean
	metric: CustomMetric | null
	saving?: boolean
}>()

const emit = defineEmits<{
	close: []
	update: [metric: CustomMetric]
}>()

const methods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'] as const
const valueTypes = ['number', 'string', 'boolean'] as const
const chartTypes = ['line', 'bar', 'pie'] as const
const timeframes = ['Last hour', 'Last 12 hours', 'Last 24 hours', 'Last 7 days', 'Last 30 days', 'Last 90 days', 'Last 365 days', 'This month', 'This year'] as const
const aggregationTypes: Record<string, string[]> = {
	number: ['count', 'sum', 'avg', 'min', 'max', 'p95', 'unique'],
	string: ['count', 'unique'],
	boolean: ['count', 'p50'],
}
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

const idPrefix = computed(() => {
	if (!props.metric) return 'update-metric'
	return `update-metric-${props.metric.name}-${props.metric.path}`.replace(/[^a-zA-Z0-9-]/g, '-')
})

const metricName = ref('')
const method = ref('POST')
const path = ref('')
const pathError = ref('')
const source = ref<MetricSource>('body')
const selector = ref('')
const aggregation = ref('avg')
const valueType = ref('number')
const timeframe = ref('Last hour')
const chartType = ref<ChartType>('line')
const applyRules = ref(false)

const syncFromMetric = (metric: CustomMetric) => {
	metricName.value = metric.name
	method.value = metric.method
	path.value = metric.path
	source.value = metric.source
	selector.value = metric.selector
	aggregation.value = metric.aggregation
	valueType.value = metric.valueType
	timeframe.value = metric.timeframe
	chartType.value = metric.chartType
	applyRules.value = metric.applyRules
	pathError.value = ''
}

watch(() => props.metric, (metric) => {
	if (metric) syncFromMetric(metric)
}, { immediate: true, deep: true })

watch(valueType, (t) => {
	const allowed = aggregationTypes[t]!
	if (!allowed.includes(aggregation.value)) {
		aggregation.value = allowed[0]!
	}
})

watch(source, (src) => {
	if (src === 'query') {
		applyRules.value = false
	}
})

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

const canSave = computed(() =>
	metricName.value.trim() !== '' && path.value.trim() !== '' && pathError.value === '' && selector.value.trim() !== '',
)

const submit = () => {
	if (!canSave.value || !props.metric) return
	emit('update', {
		...props.metric,
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
	})
}
</script>
