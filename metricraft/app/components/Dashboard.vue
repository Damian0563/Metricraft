<template>
	<div class="min-h-screen pl-20">
		<DashboardNav />
		<div class="min-w-0 mt-4">
			<GraphGrid v-if="!settings" :metrics="derivedMetrics" @load="emit('load')" />
			<Settings v-if="settings" :logRetention="logRetention" :derivedMetrics="derivedMetrics"
				@customize-view="displayView = true" @load="emit('load')" @update-metrics="emit('updateMetrics', $event)"
				@change-retention=" emit('changeRetention',
					$event)" />
			<DisplayViewCustomizer v-if="displayView" :metrics="derivedMetrics" @close="displayView = false" />
		</div>
	</div>
</template>

<script setup lang="ts">
const props = defineProps<{
	logRetention: number;
	derivedMetrics: Record<string, { enabled: boolean, timeframe: string }>;
}>();
const emit = defineEmits<{
	load: [value: void];
	updateMetrics: [value: { name: string; enabled: boolean; timeframe?: string }[]];
	changeRetention: [value: Number];
}>();
const derivedMetrics = toRef(props, 'derivedMetrics');
const logRetention = ref(props.logRetention);
const displayView = ref(false);
const route = useRoute()
const settings = computed(() => 'settings' in route.query)
watch(() => props.logRetention, (val) => logRetention.value = val);
</script>
