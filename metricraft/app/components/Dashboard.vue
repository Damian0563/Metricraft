<template>
	<div class="min-h-screen pl-20">
		<DashboardNav />
		<div class="min-w-0 mt-4">
			<GraphGrid v-if="!settings" :layout="layout" :metrics="derivedMetrics" :showView="displayView"
				@load="emit('load')" @close="displayView = false" @resetLayout="displayView = false" />
			<Settings v-if="settings" :derivedMetrics="derivedMetrics"
				@customize-view="(displayView = true, navigateTo('/dashboard'))" @load="emit('load')"
				@update-metrics="emit('updateMetrics', $event)" />
		</div>
	</div>
</template>

<script setup lang="ts">
import { useDisplayView } from "~/composables/helpers"
const props = withDefaults(defineProps<{
	derivedMetrics: Record<string, { enabled: boolean, timeframe: string }>;
	layout?: { name: string, span: number, height: number, custom: boolean }[];
}>(), {
	layout: () => [],
});
const emit = defineEmits<{
	load: [value: void];
	updateMetrics: [value: { name: string; enabled: boolean; timeframe?: string }[]];
}>();
const derivedMetrics = toRef(props, 'derivedMetrics');
const layout = toRef(props, 'layout');
const displayView = useDisplayView();
const route = useRoute()
const settings = computed(() => 'settings' in route.query)
watch(settings, (on) => { if (!on) displayView.value = false })
watch(() => props.logRetention, (val) => logRetention.value = val);
</script>
