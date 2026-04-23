<template>
	<NuxtLayout>
		<DashboardNav :appName="appName" :realtimeEnabled="realtimeEnabled" @realtimeToggle=handleRealtimeToggle
			@customizeView=handleCustomizeView />
		<GraphGrid />
	</NuxtLayout>
</template>

<script setup lang="ts">
import { handleMessage } from "@/ws/visitors"
import { toggleRealtime } from "@/calls/settings"
const props = defineProps<{
	appName: string;
	realtimeEnabled: boolean;
}>();
const appName = ref(props.appName);
const realtimeEnabled = ref(props.realtimeEnabled);
watch(() => props.appName, (val) => appName.value = val);
watch(() => props.realtimeEnabled, (val) => realtimeEnabled.value = val);
const emit = defineEmits<{
	load: void;
}>();
let ws: WebSocket | null = null;
const connectWebSocket = () => {
	ws = new WebSocket(`ws://localhost:8080/ws/visitors`);
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
