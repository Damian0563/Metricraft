<template>
	<NuxtLayout>
		<DashboardNav :appName="appName" @settings=handleSettings />
		<GraphGrid v-if="!settings" />
		<Settings v-if="settings" :appName="appName" :realtimeEnabled="realtimeEnabled"
			@realtime-toggle="handleRealtimeToggle" @customize-view="handleCustomizeView" @load="emit('load')"
			@back="handleSettings" />
	</NuxtLayout>
</template>

<script setup lang="ts">
import { handleMessage } from "@/ws/visitors"
import { toggleRealtime } from "@/calls/settings"
import type { config } from "@/composables/types"
const props = defineProps<{
	appName: string;
	realtimeEnabled: boolean;
}>();
const appName = ref(props.appName);
const realtimeEnabled = ref(props.realtimeEnabled);
const settings = ref(false);
watch(() => props.appName, (val) => appName.value = val);
watch(() => props.realtimeEnabled, (val) => realtimeEnabled.value = val);
const emit = defineEmits<{
	load: void;
}>();
let ws: WebSocket | null = null;
const connectWebSocket = () => {
	const config: config = useBackendUrl()
	if (!config.wsshost) console.error("No websocket host")
	ws = new WebSocket(`${config.wsshost}/ws/visitors`);
	ws.onopen = () => {
		console.log("Connected to websocket");
	};
	ws.onmessage = (event: MessageEvent) => {
		console.log(event);
		handleMessage(event);
	};
};

const handleRealtimeToggle = (val: boolean) => {
	toggleRealtime(val)
	if (val) {
		connectWebSocket();
	} else {
		console.log("Disconnected from websocket");
		if (ws) ws.close();
	}
};

const handleSettings = () => {
	settings.value = !settings.value;
}

const handleCustomizeView = (val: boolean) => {
	console.log(val);
};

onMounted(() => {
	if (realtimeEnabled.value) {
		connectWebSocket();
	}
});

onBeforeUnmount(() => {
	if (ws) ws.close();
});
</script>
