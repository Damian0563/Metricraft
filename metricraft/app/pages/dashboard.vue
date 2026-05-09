<template>
	<div>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Spinner :loading="loading || localLoading" />
		<Dashboard :realtimeEnabled="realtimeEnabled" :logRetention="logRetention" :derivedMetrics="derivedMetrics"
			@load="handleLoad" />
	</div>
</template>

<script setup lang="ts">
import type { dashboardInitPayload } from '@/composables/types';
import { getDashboard } from "~/calls/dashboard";
const localLoading = ref(false);
const errorMessage = ref("");
const appName = useState<string>('appName', () => "");
const realtimeEnabled = ref(false);
const logRetention = ref(30);
const derivedMetrics = ref(new Map<string, boolean>());
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
		derivedMetrics.value = new Map(Object.entries(raw))
	} else {
		errorMessage.value = newVal.error
		timeout.value = setTimeout(() => {
			navigateTo("/")
		}, 5200)
	}
})
const handleLoad = () => {
	localLoading.value = !localLoading.value
};
initialize(payload.value)
localLoading.value = loading.value
if (error.value) {
	errorMessage.value = "Something went wrong, Check your internet connection and try again."
}
onBeforeUnmount(() => {
	clearTimeout(timeout.value)
})
</script>
