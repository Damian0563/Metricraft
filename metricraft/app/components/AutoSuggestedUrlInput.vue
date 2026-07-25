<template>
	<div>
		<div class="flex items-baseline justify-between gap-3 mb-2">
			<label :for="inputId" class="text-sm font-medium text-gray-700">{{ label }}</label>
			<p v-if="patternError" class="text-xs text-red-500 text-right">{{ patternError }}</p>
		</div>
		<div class="relative">
			<input :id="inputId" v-model="pattern" type="text" :placeholder="placeholder" autocomplete="off"
				spellcheck="false" @focus="onBaseFocus" @blur="onBaseBlur" @input="showSuggestedUrls" @keyup="showSuggestedUrls"
				@keyup.enter="emit('submitEvent')"
				class="w-full px-4 py-2 rounded-lg border border-gray-200 text-gray-800 font-mono text-sm focus:outline-none focus:border-[#00F376] transition-colors" />
			<ul v-if="inputActive && suggestedUrls.length > 0"
				class="absolute left-0 right-0 top-full z-20 mt-1 max-h-56 overflow-y-auto rounded-lg border border-gray-200 bg-white shadow-xl">
				<li v-for="url in suggestedUrls" :key="url">
					<button type="button" @mousedown.prevent="selectSuggestion(url)"
						class="block w-full truncate px-4 py-2 text-left font-mono text-sm text-gray-700 hover:bg-[#00F376]/10 hover:text-[#00B35C] transition-colors cursor-pointer">
						{{ url }}
					</button>
				</li>
			</ul>
		</div>
		<p class="mt-2 text-xs text-gray-500">{{ helpText }}</p>
	</div>
</template>

<script setup lang="ts">
import { getUrls } from '@/calls/dashboard'
const pattern = defineModel<string>({ default: '' })

const props = withDefaults(defineProps<{
	groupingMode?: boolean
	existingRules?: string[]
	label?: string
	helpText?: string
	inputId?: string
	placeholder?: string
}>(), {
	label: 'Base path',
	inputId: 'rule-pattern',
	placeholder: 'https://api.service/users/id',
})

const helpText = computed(() => {
	if (props.helpText) return props.helpText
	return props.groupingMode
		? 'Routes are matched by prefix. Deeper segments become the collapsed tail.'
		: 'Routes are matched by prefix. Deeper segments become blacklisted too.'
})

const emit = defineEmits<{
	submitEvent: []
	patternError: [value: string]
	fetchedUrls: [value: string[]]
}>()

const inputActive = ref(false)
const suggestedUrls = ref<string[]>([])
const urls = useState<string[]>('urls', () => [])
if (urls.value.length === 0) {
	const { data: fetchedUrls } = await useAsyncData('dashboardUrls', () => getUrls())
	if (fetchedUrls.value) urls.value = fetchedUrls.value
}
emit('fetchedUrls', urls.value)


const urlHasPrefix = (path: string, url: string): boolean => {
	if (!path.startsWith(url)) return false
	if (path === url) return true
	return path.slice(url.length).startsWith('/')
}

const pathnameOf = (value: string): string => {
	try {
		return new URL(value).pathname
	} catch {
		const trimmed = value.trim()
		return trimmed ? `/${trimmed}` : trimmed
	}
}

const urlMatchesQuery = (query: string, url: string): boolean => {
	if (!query) return true
	if (url.startsWith(query)) return true
	return pathnameOf(url).startsWith(pathnameOf(query))
}

const normalizedBase = computed(() => pattern.value.trim())

const prefixExists = (path: string): boolean =>
	urls.value.some((url) => urlHasPrefix(path, url))

const patternError = computed(() => {
	const base = normalizedBase.value
	if (!base) return ''
	if (base.includes('http') && base.split('/').length <= 3) {
		return `Add at least one path segment, e.g. ${base}/more.`
	}
	if (base.at(-1) === '/') {
		return `Remove the trailing slash, e.g. ${base.slice(0, -1)}.`
	}
	if (props.existingRules?.includes(base)) {
		return 'A rule for this path already exists.'
	}
	if (!prefixExists(base)) {
		return 'The base path is not in the list of URLs.'
	}
	return ''
})

watch(patternError, (message) => emit('patternError', message), { immediate: true })

const onBaseFocus = () => {
	inputActive.value = true
	showSuggestedUrls()
}

const onBaseBlur = () => {
	inputActive.value = false
}

const selectSuggestion = (url: string) => {
	pattern.value = url
	inputActive.value = false
	showSuggestedUrls()
}

const showSuggestedUrls = () => {
	const query = pattern.value.trim()
	suggestedUrls.value = urls.value.filter((url) => urlMatchesQuery(query, url)).slice(0, 6)
}
</script>
