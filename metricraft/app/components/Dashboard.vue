<template>
	<div>
		<DashboardNav @settings=handleSettings />
		<GraphGrid v-if="!settings" />
		<Settings v-if="settings" :realtimeEnabled="realtimeEnabled" :logRetention="logRetention"
			:derivedMetrics="derivedMetrics" @realtime-toggle="handleRealtimeToggle" @customize-view="handleCustomizeView"
			@load="emit('load')" @update-metrics="emit('updateMetrics', $event)" @change-retention=" emit('changeRetention',
				$event)" />
	</div>
</template>

<script setup lang="ts">
import { handleMessage } from "@/ws/visitors"
import { toggleRealtime } from "@/calls/settings"
import type { config } from "@/composables/types"
const props = defineProps<{
	realtimeEnabled: boolean;
	logRetention: number;
	derivedMetrics: Record<string, boolean>;
}>();
const emit = defineEmits<{
	load: [value: void];
	updateMetrics: [value: { name: string; enabled: boolean }[]];
	changeRetention: [value: Number];
	changeRealtime: [value: boolean];
}>();
const derivedMetrics = toRef(props, 'derivedMetrics');
const realtimeEnabled = ref(props.realtimeEnabled);
const logRetention = ref(props.logRetention);
const route = useRoute()
const settings = computed(() => 'settings' in route.query)
watch(() => props.realtimeEnabled, (val) => realtimeEnabled.value = val);
watch(() => props.logRetention, (val) => logRetention.value = val);
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
	emit('changeRealtime', val)
	toggleRealtime(val)
	if (val) {
		connectWebSocket();
	} else {
		console.log("Disconnected from websocket");
		if (ws) ws.close();
	}
};

const handleSettings = () => {
	navigateTo('/dashboard?settings')
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
