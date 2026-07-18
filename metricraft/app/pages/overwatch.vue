<template>
	<div class="min-h-screen pl-20">
		<DashboardNav />
		<ClientOnly>
			<Spinner :loading="loading" />
		</ClientOnly>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<div class="min-w-0 px-8 py-2">
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
								<label for="metric-name" class="block text-sm font-medium text-gray-700 mb-2">Metric name</label>
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
									<label for="metric-path" class="block text-sm font-medium text-gray-700 mb-2">Endpoint path</label>
									<input id="metric-path" v-model="path" type="text" placeholder="/api/v1/checkout"
										class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors" />
									<p class="mt-2 text-xs text-gray-500">Overwatch inspects requests matching this route.</p>
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
										<option v-for="agg in aggregationTypes[valueType]" :key="agg" :value="agg">{{
											agg }}</option>
									</select>
								</div>
								<div class="flex-1 min-w-0">
									<label for="metric-timeframe" class="block text-sm font-medium text-gray-700 mb-2">Timeframe</label>
									<select id="metric-timeframe" v-model="timeframe"
										class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 focus:outline-none focus:border-[#00F376] transition-colors cursor-pointer">
										<option v-for="tf in timeframes" :key="tf" :value="tf">{{ tf }}</option>
									</select>
								</div>
							</div>

							<div class="pt-6 border-t border-gray-100 flex justify-end gap-3">
								<button type="button" @click="resetForm"
									class="px-6 py-3 bg-gray-100 text-gray-700 font-semibold rounded-lg hover:bg-gray-200 transition-colors cursor-pointer">
									Reset
								</button>
								<button type="button" @click="addMetric" :disabled="!canSave"
									class="px-8 py-3 bg-[#00F376] text-gray-900 font-bold rounded-lg hover:bg-[#00D96A] transition-all shadow-lg disabled:opacity-50 disabled:cursor-not-allowed uppercase tracking-wider text-sm cursor-pointer">
									Create metric
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
										:class="line.target ? 'bg-[#00F376]/20 ring-1 ring-[#00F376]/50 px-0.5 text-emerald-300' : ''">{{ line.text }}<span
											v-if="line.target" class="text-[#00F376]/70"> ← inspected</span></div>
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
									<span class="font-medium text-gray-700">{{ timeframe.toLowerCase() }}</span>
								</div>
							</div>
						</div>
					</section>
					<section class="flex flex-col bg-gray-50/60">
						<div class="px-6 py-5 border-b border-gray-100">
							<div class="flex items-center justify-between gap-3">
								<h2 class="text-sm font-semibold text-gray-800">Configured metrics</h2>
								<span class="shrink-0 text-xs font-medium text-gray-500">{{ metrics.length }} defined</span>
							</div>
						</div>
					</section>
				</div>
			</div>
		</div>
	</div>
</template>


<script setup lang="ts">
type MetricSource = 'body' | 'header' | 'query'
interface CustomMetric {
	name: string
	method: string
	path: string
	source: MetricSource
	selector: string
	aggregation: string
	timeframe: string
	valueType: string
}
const loading = ref(false)
const errorMessage = ref('')
const metricName = ref('')
const method = ref<string>('POST')
const path = ref('')
const source = ref<MetricSource>('body')
const selector = ref('')
const aggregation = ref<string>('avg')
const metrics = ref<CustomMetric[]>([])
const sourceLabel = computed(() => sources.find(s => s.id === source.value)?.label ?? '')
const methods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'] as const
const valueType: Ref<string> = ref('number')
const aggregationTypes: { [key: string]: string[] } = {
	'number': ['count', 'sum', 'avg', 'min', 'max', 'p95', 'unique'],
	'string': ['count', 'unique'],
	'boolean': ['count', 'p50'],
}
const valueTypes = ['number', 'string', 'boolean'] as const
watch(valueType, (t) => {
	const allowed = aggregationTypes[t]!
	if (!allowed.includes(aggregation.value)) {
		aggregation.value = allowed[0]!
	}
})
const timeframes = ['Last hour', 'Last 12 hours', 'Last 24 hours', 'Last 7 days', 'Last 30 days', 'Last 90 days', 'Last 365 days', 'This month', 'This year'] as const
const timeframe = ref<string>('Last hour')
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
interface BodyLine { text: string; target: boolean }
const TARGET_TOKEN = '__OVERWATCH_TARGET__'
const buildBodyObject = (segments: string[]): unknown => {
	const [head, ...rest] = segments
	if (head === undefined) return TARGET_TOKEN
	if (/^\d+$/.test(head)) return [buildBodyObject(rest)]
	return { [head]: buildBodyObject(rest) }
}

const bodyLines = computed<BodyLine[]>(() => {
	const selecting = source.value === 'body'
	const bodyObject = selecting
		? buildBodyObject(selector.value.trim().replace(/\[(\d+)\]/g, '.$1').split('.').map(s => s.trim()).filter(Boolean))
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
	metricName.value.trim() !== '' && path.value.trim() !== '' && selector.value.trim() !== '')
const addMetric = () => {
	if (!canSave.value) return
	metrics.value.push({
		name: metricName.value.trim(),
		method: method.value,
		path: path.value.trim(),
		source: source.value,
		selector: selector.value.trim(),
		timeframe: timeframe.value,
		aggregation: aggregation.value,
		valueType: valueType.value,
	})
	resetForm()
}

const resetForm = () => {
	metricName.value = ''
	method.value = 'POST'
	path.value = ''
	source.value = 'body'
	selector.value = ''
	aggregation.value = 'avg'
	valueType.value = 'number'
}
</script>
