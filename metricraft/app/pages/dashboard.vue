<template>
	<NuxtLayout>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Spinner :loading="loading" />
		<Dashboard :realtimeEnabled="realtimeEnabled" :logRetention="logRetention" :derivedMetrics="derivedMetrics"
			@load="handleLoad" />
	</NuxtLayout>
</template>

<script setup lang="ts">
import { getCookie, updateCookie } from "@/composables/helpers";
import type { dashboardInitPayload } from '@/composables/types';
import { getDashboard } from "~/calls/dashboard";
const cookie = ref("");
const loading = ref(true);
const errorMessage = ref("");
const appName = useState<string>('appName', () => "");
const realtimeEnabled = ref(false);
const logRetention = ref(30);
const derivedMetrics = ref(new Map<string, boolean>());
const init = async () => {
	loading.value = true;
	try {
		const data: dashboardInitPayload = await getDashboard(cookie.value);
		if (data.error) {
			errorMessage.value = "Error loading dashboard, session expired.";
			navigateTo("/");
		} else {
			appName.value = data.appName;
			realtimeEnabled.value = data.settings.realtime;
			logRetention.value = data.settings.retention;
			const raw = data.settings.enabled as Record<string, boolean>;
			derivedMetrics.value = new Map(Object.entries(raw));
			cookie.value = data.signedSecret;
			updateCookie(cookie.value);
		}
	} catch (e) {
		console.log(e);
		errorMessage.value = "Error loading dashboard, session expired.";
		navigateTo("/");
	} finally {
		loading.value = false;
	}
};

const handleLoad = () => {
	loading.value = !loading.value;
};

onMounted(() => {
	cookie.value = getCookie("session-token");
	init();
});
</script>
