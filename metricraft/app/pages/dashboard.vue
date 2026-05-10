<template>
	<div>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Spinner :loading="loading || localLoading" />
		<Dashboard :realtimeEnabled="realtimeEnabled" :logRetention="logRetention" :derivedMetrics="derivedMetrics"
			@load="handleLoad" @updateMetrics="handleUpdateMetrics" />
	</div>
</template>

<script setup lang="ts">
import type { dashboardInitPayload } from '@/composables/types';
import { getDashboard, fetchMetric } from "~/calls/dashboard";
const localLoading = ref(false);
const errorMessage = ref("");
const appName = useState<string>('appName', () => "");
const derivedMetrics = ref<Record<string, boolean>>({})
const realtimeEnabled = ref(false);
const logRetention = ref(30);
const timeout = ref(0)
const { data: payload, pending: loading, error } = await useAsyncData<dashboardInitPayload>('dashboard', () => getDashboard())
const initialize = ((newVal: dashboardInitPayload | undefined) => {
	if (!newVal || newVal === undefined) {
		return
	}
	if (newVal && newVal.error === '') {
		appName.value = newVal.appName
		realtimeEnabled.value = newVal.settings.realtime
		logRetention.value = newVal.settings.retention
		const raw = newVal.settings.enabled as Record<string, boolean>
		derivedMetrics.value = raw
	} else {
		errorMessage.value = newVal.error
		if (import.meta.client) {
			timeout.value = setTimeout(() => {
				navigateTo("/")
			}, 5200)
		}
	}
})

watch(() => derivedMetrics.value, (val: Record<string, boolean>) => {
	if (!val || val === undefined) return
	Object.entries(val).forEach(([name, enabled]) => {
		if (enabled) {
			void fetchMetric(name)
		}
	})
})
})
const handleLoad = () => {
	localLoading.value = !localLoading.value
};
const handleUpdateMetrics = (changes: { name: string; enabled: boolean }[]) => {
	derivedMetrics.value = Object.fromEntries(changes.map(c => [c.name, c.enabled]))
};
watch(() => payload.value, initialize, { immediate: true })
localLoading.value = loading.value
if (error.value) {
	errorMessage.value = "Something went wrong, Check your internet connection and try again."
}
onBeforeUnmount(() => {
	clearTimeout(timeout.value)
})
</script>
