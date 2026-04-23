<template>
	<NuxtLayout>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Spinner :loading="loading" />
		<Dashboard :appName="appName" @load="handleLoad" />
	</NuxtLayout>
</template>

<script setup lang="ts">
import { getCookie, updateCookie } from "@/composables/helpers";
import type { dashboardInitPayload } from '@/composables/types';
import { getDashboard } from "~/calls/dashboard";
const cookie = ref("");
const loading = ref(true);
const errorMessage = ref("");
const appName = ref("");

const init = async () => {
	loading.value = true;
	try {
		const data: dashboardInitPayload = await getDashboard(cookie.value);
		console.log("raw response:", JSON.stringify(data));
		if (data.error) {
			errorMessage.value = "Error loading dashboard, session expired.";
			navigateTo("/");
		} else {
			appName.value = data.appName;
			cookie.value = data.signedSecret;
			updateCookie(cookie.value);
		}
	} catch (e) {
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
