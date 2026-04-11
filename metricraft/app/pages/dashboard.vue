<template>
	<NuxtLayout>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Spinner :loading="loading" />
	</NuxtLayout>
</template>

<script setup lang="ts">
import { getCookie, updateCookie } from "@/composables/helpers";
import type { dashboardInitPayload } from '@/composables/types';
import { getDashboard } from "~/calls/dashboard";
import { handleMessage } from "@/ws/visitors"

const cookie = ref("");
const loading = ref(true);
const errorMessage = ref("");
const appName = ref("");
let ws: WebSocket | null = null;

const init = async () => {
	loading.value = true;
	try {
		const data: dashboardInitPayload = await getDashboard(cookie.value);
		if (data.error) {
			errorMessage.value = data.error;
		} else {
			appName.value = data.appName;
			cookie.value = data.signedSecret;
			updateCookie(cookie.value);
		}
	} catch (e) {
		errorMessage.value = e as string;
	} finally {
		loading.value = false;
	}
};

onMounted(() => {
	cookie.value = getCookie("session-token");
	init();
	ws = new WebSocket(`ws://localhost:8080/ws/visitors`);
	ws.onopen = () => {
		console.log("Connected to websocket");
	};
	ws.onmessage = (event: MessageEvent) => {
		console.log(event);
		handleMessage(event);
	};
});

onBeforeUnmount(() => {
	if (ws) ws.close();
});
</script>
