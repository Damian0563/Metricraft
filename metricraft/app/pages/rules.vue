<template>
	<div>
		<ClientOnly>
			<Spinner :loading="saving" />
		</ClientOnly>
		<Notice :message="errorMessage ?? ''" @close="errorMessage = null" />
		<Popup :message="cleanMessage" @close="cleanMessage = ''" />
		<div class="relative flex items-center justify-center mb-6">
			<h1 class="text-3xl font-bold text-center" style="color: #00F376;">Rules</h1>
		</div>
		<div class="justify-between mb-2">
			<button type="button" @click="router.replace({ query: { type: 'grouping' } })"
				class="px-2 py-2 text-gray-700 font-semibold rounded-lg hover:opacity-80 transition-colors cursor-pointer"
				:class="groupingMode ? 'bg-[#00F376]' : 'bg-gray-100'">
				Grouping
			</button>
			<button type="button" @click="router.replace({ query: { type: 'blacklisting' } })"
				class="px-2 py-2 mx-2 text-gray-700 font-semibold rounded-lg hover:opacity-80 transition-colors cursor-pointer"
				:class="!groupingMode ? 'bg-[#00F376]' : 'bg-gray-100'">
				Blacklisting
			</button>
		</div>
		<div class="max-w-8xl mx-auto grid gap-4 lg:grid-cols-[minmax(0,1fr)_40rem] lg:items-stretch">
			<div class="flex flex-col gap-4 min-w-0">
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<h2 class="text-xl font-semibold text-gray-800 mb-4">
						What are {{ groupingMode ? 'Grouping' : 'Blacklisting' }} rules?
					</h2>
					<div class="text-sm text-gray-600 leading-relaxed mb-4">
						<div v-if="groupingMode">
							Rules group routes that share a common path but differ only by a trailing
							identifier. Instead of charting <span class="font-mono text-gray-800">/users/id/1</span>,
							<span class="font-mono text-gray-800">/users/id/2</span> and
							<span class="font-mono text-gray-800">/users/id/3</span> as separate endpoints, a rule
							folds them into a single <span class="font-mono text-gray-800">/users/more/</span> line on
							your dashboards. Rules can be also applied to query parameters in the case <span
								class="font-mono text-gray-800">/users/id?value=1</span> and <span
								class="font-mono text-gray-800">/users/id?value=2</span> are collapsed into a single
							<span class="font-mono text-gray-800">/users/id</span> line.
						</div>
						<div v-else>
							Blacklist rules omit a given explicit endpoint path or implicit endpoints that have path parameters and
							query parameters from your dashboards. The former is useful for hiding non-user facing endpoints used
							internally i.e for telemetry or healthchecks.
							Implicit blacklisting allows you to granualarly hide endpoints that match a given pattern. In practice,
							if <span class="font-mono text-gray-800">/users/id</span>
							is blacklisted then <span class="font-mono text-gray-800">/users/id/1</span> and subsequent child routes
							are not included in your dashboard.
						</div>
					</div>
					<div class="text-sm text-gray-600 leading-relaxed">
						<div v-if="groupingMode">
							Point a rule at a base path and every deeper route will collapses into it. Your metrics stay
							readable even when a route explodes into thousands of unique URLs. When adding a rule remember about
							supplying a valid path, the one which is the target of the request,<strong> in particular remember about
								the
								full prefix, in the above examples it could be https://api.service. </strong>
						</div>
						<div v-else>
							Point a rule at a base path that will be removed from chart composition. When adding a rule remember about
							supplying a valid path, the one which is the target of the requests,<strong> in particular remember about
								the
								full prefix, in the above examples it could be https://api.service. </strong>
						</div>
					</div>
				</div>

				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<div class="flex items-start gap-3 mb-6">
						<div
							class="h-10 w-10 shrink-0 rounded-full bg-[#00F376]/10 flex items-center justify-center text-[#00B35C]">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"
								aria-hidden="true">
								<path fill-rule="evenodd"
									d="M3 4a1 1 0 011-1h3a1 1 0 010 2H5v10h10v-2a1 1 0 112 0v3a1 1 0 01-1 1H4a1 1 0 01-1-1V4zm9.293-.707a1 1 0 011.414 0l3 3a1 1 0 010 1.414l-6 6A1 1 0 0110 14H8a1 1 0 01-1-1v-2a1 1 0 01.293-.707l6-6z"
									clip-rule="evenodd" />
							</svg>
						</div>
						<div class="min-w-0">
							<h2 class="text-xl font-semibold text-gray-800">Create a {{ groupingMode ? 'grouping' : 'blacklisting'
							}} rule</h2>
							<p class="text-sm text-gray-500 mt-1">
								{{ groupingMode ? 'Enter the base path to keep. Everything below it collapses into one endpoint.'
									: 'Enter the base path to hide. Everything below it will be omitted from your dashboards.'
								}}
							</p>
						</div>
					</div>

					<div class="space-y-6">
						<AutoSuggestedUrlInput v-model="pattern" :grouping-mode="groupingMode"
							:existing-rules="rulesList.map((rule) => rule.rule)" @submit-event="addRuleEntry"
							@pattern-error="patternError = $event" @fetched-urls="urls = $event" />
						<div v-if="groupingMode">
							<p class="text-sm font-medium text-gray-700 mb-2">How it collapses</p>
							<div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
								<div class="grid items-center gap-3 sm:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)]">
									<div class="space-y-1.5">
										<p class="text-[0.65rem] font-semibold uppercase tracking-wider text-gray-400 mb-1"
											v-if="samplePaths.length > 1">
											{{ samplePaths.length }} distinct routes
										</p>
										<p class="text-[0.65rem] font-semibold uppercase tracking-wider text-orange-500 mb-1"
											v-else-if="samplePaths.length === 1">
											Only one route falls under this rule in traffic observed thus far. Before more routes are matched
											this route will serve merely as a rename of the grouped route.
										</p>
										<p class="text-[0.65rem] font-semibold uppercase tracking-wider text-red-500 mb-1" v-else>
											No routes fall under this rule, thus far no such traffic has been observed. If this is intended
											you can add such rule anyway- it will be populated after traffic starts flowing.
										</p>
										<div v-for="path in samplePaths" :key="path.full"
											class="whitespace-nowrap overflow-x-auto rounded-md bg-white border border-gray-200 px-2.5 py-1.5 font-mono text-xs text-gray-500">
											<span class="text-gray-700">{{ path.base }}</span><span class="text-gray-300"></span><span
												class="rounded bg-gray-100 px-1 text-gray-600">{{ path.tail }}</span>
										</div>
									</div>
									<div class="flex items-center justify-center text-[#00B35C]" aria-hidden="true">
										<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 rotate-90 sm:rotate-0" viewBox="0 0 20 20"
											fill="currentColor">
											<path fill-rule="evenodd"
												d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z"
												clip-rule="evenodd" />
										</svg>
									</div>

									<div class="space-y-1.5">
										<p class="text-[0.65rem] font-semibold uppercase tracking-wider text-gray-400 mb-1">
											1 grouped endpoint
										</p>
										<div
											class="whitespace-nowrap overflow-x-auto rounded-md bg-white px-2.5 py-1.5 font-mono text-xs ring-1 ring-[#00F376]/50 border border-[#00F376]">
											<span class="text-gray-800 font-semibold">{{ normalizedBase || 'https://api.service/users/id'
											}}</span>
										</div>
										<p class="text-xs text-gray-500 pt-1">
											Charted once, with the request counts of every child summed together.
										</p>
									</div>
								</div>
							</div>
						</div>

						<div class="pt-6 border-t border-gray-100 flex justify-end gap-3">
							<button type="button" @click="resetForm"
								class="px-6 py-3 bg-gray-100 text-gray-700 font-semibold rounded-lg hover:bg-gray-200 transition-colors cursor-pointer">
								Reset
							</button>
							<button type="button" @click="addRuleEntry" :disabled="!canSave || saving"
								class="px-8 py-3 bg-[#00F376] text-gray-900 font-bold rounded-lg hover:bg-[#00D96A] transition-all shadow-lg disabled:opacity-50 disabled:cursor-not-allowed uppercase tracking-wider text-sm cursor-pointer">
								Create rule
							</button>
						</div>
					</div>
				</div>
			</div>

			<div
				class="bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden lg:sticky lg:top-8 flex flex-col">
				<div class="shrink-0 px-6 py-5 border-b border-gray-100">
					<div class="flex items-center justify-between gap-3">
						<h2 class="text-xl font-semibold text-gray-800">Active rules</h2>
						<span class="shrink-0 text-xs font-medium text-gray-500">
							{{rulesList.filter(rule => groupingMode ? rule.mode === 'grouping' : rule.mode ===
								'blacklisting').length}} rule{{
								rulesList.filter(rule => groupingMode ? rule.mode === 'grouping' : rule.mode === 'blacklisting').length
									=== 1 ? '' :
									's'}}
						</span>
					</div>
					<p class="text-sm text-gray-500 mt-1">
						Every rule below is applied to your traffic before it reaches your charts.
					</p>
				</div>

				<div v-if="rulesList.length > 0" class="min-h-0 flex-1 overflow-y-auto">
					<ClientOnly>
						<AnimatePresence>
							<motion.div v-for="(entry, index) in rulesList" layout :key="entry.rule"
								:initial="{ opacity: 0, height: 0 }" :animate="{ opacity: 1, height: 'auto' }"
								:exit="{ opacity: 0, height: 0, x: 32, transition: { duration: 0.22, ease: 'easeIn' } }"
								:transition="{ type: 'spring', stiffness: 480, damping: 34, delay: Math.min(index * 0.045, 0.27) }">
								<div v-if="groupingMode ? entry.mode === 'grouping' : entry.mode === 'blacklisting'"
									class="flex items-center gap-3 px-6 py-4 border-b border-gray-100 last:border-b-0 hover:bg-gray-50 transition-colors">
									<div class="min-w-0 flex-1">
										<p class="font-mono text-sm text-gray-900 break-all leading-snug">
											{{ entry.rule }}<span class="text-gray-300">/</span><span
												class="text-[#00B35C] font-semibold">*</span>
										</p>
										<p class="mt-1 text-xs"
											:class="entry.matches.length > 0 ? 'text-gray-500' : 'text-gray-400 italic'">
											{{ entry.matches.length > 5
												? `Collapses ${entry.matches.slice(0, 5)} routes and ${entry.matches.length - 5} more`
												: entry.matches.length > 0 ? `Collapses ${entry.matches} route${entry.matches.length === 1 ? '' :
													's'}` : 'Waiting for matching traffic' }}
										</p>
									</div>
									<button type="button" @click="removeRule(entry)" :disabled="saving"
										class="shrink-0 px-4 py-2 text-sm font-semibold rounded-lg bg-red-50 text-red-600 hover:bg-red-100 transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed">
										Delete
									</button>
								</div>
							</motion.div>
						</AnimatePresence>
					</ClientOnly>
				</div>

				<div v-else class="min-h-0 flex-1 px-6 py-10 text-center">
					<div class="mx-auto mb-3 h-10 w-10 rounded-full bg-gray-100 flex items-center justify-center text-gray-400">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd"
								d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z"
								clip-rule="evenodd" />
						</svg>
					</div>
					<p class="text-sm text-gray-500">No rules yet.</p>
					<p class="text-xs text-gray-400 mt-1">Create a rule to start grouping routes.</p>
				</div>
			</div>
		</div>
	</div>
</template>


<script setup lang="ts">
definePageMeta({
	layout: 'dashboard',
})
import { getRules, addRule, deleteRule } from '@/calls/rules'
import { motion, AnimatePresence } from "motion-v"
import { parseApiError } from '@/composables/helpers'
import type { Rule } from '@/composables/types/additional'

const router = useRouter()
const groupingMode = computed(() => {
	return router.currentRoute.value.query.type === 'grouping'
})
const errorMessage = ref<string | null>(null)
const cleanMessage = ref<string>('')
const saving = ref(false)
const pattern = ref('')
const patternError = ref('')
const urls = ref<string[]>([])
const { data: fetchedRules, error: fetchError } = await useAsyncData<Rule[]>('rules', () => getRules(), { default: () => [] })
const rulesList = ref<Rule[]>([])
watch(fetchedRules, (rules) => {
	if (rules) rulesList.value = [...rules]
}, { immediate: true })
if (fetchError.value) {
	errorMessage.value = 'Failed to load rules.'
}

const normalizedBase = computed(() => {
	const trimmed = pattern.value.trim()
	if (!trimmed) return ''
	return trimmed
})

const canSave = computed(() =>
	normalizedBase.value.length > 1 && patternError.value === '')

const samplePaths = computed(() => {
	const base = normalizedBase.value || 'https://api.service/users/id'
	if (base === normalizedBase.value) {
		const actualMatched = urls.value.filter((url) => url.startsWith(base))
		return actualMatched.map((url) => ({ full: url, base, tail: url.replace(base, '') }))
	}
	return ['/1', '/2', '/87'].map((tail) => ({ full: `${base}/${tail}`, base, tail }))
})

const addRuleEntry = async () => {
	if (!canSave.value) return
	const rule = normalizedBase.value
	saving.value = true
	try {
		await addRule({ rule, mode: groupingMode.value ? 'grouping' : 'blacklisting', matches: samplePaths.value.map((path) => path.full) })
		rulesList.value.push({ rule, matches: samplePaths.value.map((path) => path.full), mode: groupingMode.value ? 'grouping' : 'blacklisting' })
		cleanMessage.value = 'Rule created successfully.'
		resetForm()
	} catch (e: unknown) {
		errorMessage.value = parseApiError(e, 'Something went wrong while saving the rule. Please try again.')
	} finally {
		saving.value = false
	}
}

const removeRule = async (rule: Rule) => {
	saving.value = true
	try {
		await deleteRule(rule)
		rulesList.value = rulesList.value.filter((entry) => entry.rule !== rule.rule)
		cleanMessage.value = 'Rule deleted successfully.'
	} catch (e: unknown) {
		errorMessage.value = parseApiError(e, 'Something went wrong while deleting the rule. Please try again.')
	} finally {
		saving.value = false
	}
}

const resetForm = () => {
	pattern.value = ''
}
</script>
