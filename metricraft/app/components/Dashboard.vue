<template>
	<NuxtLayout>
		<DashboardNav :appName="appName" @realTimeToggle=handleRealtimeToggle @customizeView=handleCustomizeView />
		<GraphGrid />
	</NuxtLayout>
</template>

<script setup lang="ts">
import { handleMessage } from "@/ws/visitors"
const props = defineProps<{
	appName: string;
}>();
const appName = ref("");
watch(() => props.appName, (val) => appName.value = val);
const emit = defineEmits(['load']);
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
	connectWebSocket();
});

onBeforeUnmount(() => {
	if (ws) ws.close();
});
</script>
