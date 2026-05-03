<template>
	<NuxtLayout>
		<DashboardNav :appName="appName" @settings=handleSettings @back="settings = false" />
		<GraphGrid v-if="!settings" />
		<Settings v-if="settings" :realtimeEnabled="realtimeEnabled" :logRetention="logRetention"
			@realtime-toggle="handleRealtimeToggle" @customize-view="handleCustomizeView" @load="emit('load')"
			@logRetention="handleLogRetention" />
	</NuxtLayout>
</template>

<script setup lang="ts">
import { handleMessage } from "@/ws/visitors"
import { toggleRealtime, changeRetention } from "@/calls/settings"
import type { config } from "@/composables/types"
const props = defineProps<{
	appName: string;
	realtimeEnabled: boolean;
	logRetention: number;
}>();
const appName = ref(props.appName);
const realtimeEnabled = ref(props.realtimeEnabled);
const logRetention = ref(props.logRetention);
const settings = ref(false);
watch(() => props.appName, (val) => appName.value = val);
watch(() => props.realtimeEnabled, (val) => realtimeEnabled.value = val);
watch(() => props.logRetention, (val) => logRetention.value = val);
const emit = defineEmits<{
	load: [value: void];
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

const handleLogRetention = (val: number) => {
	changeRetention(val);
}

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
	settings.value = true;
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
