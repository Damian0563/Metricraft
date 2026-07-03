<template>
	<div class="flex justify-between items-center">
		<p v-if="aux_header"
			class="shrink-0 text-base font-semibold tabular-nums tracking-tight text-slate-400 font-size-sm">
			{{ aux_header }}
		</p>
		<p v-if="header" class="shrink-0 text-base font-semibold tabular-nums tracking-tight text-slate-400 font-size-sm">
			{{ header }}
		</p>
	</div>
</template>

<script setup lang="ts">
const props = defineProps<{
	data: { metric: string, data: any }
}>();
const header = ref<string | null>(null);
const aux_header = ref<string | null>(null);
watch(() => props.data, (newData: any) => {
	if (newData.metric === "Throughput") {
		if (newData.data.computedThroughput === undefined || newData.data.uniqUsers === undefined) return;
		header.value = `${newData.data.computedThroughput.toFixed(4)} requests/s`
		aux_header.value = `${newData.data.uniqUsers} unique users`
	}
}, { immediate: true })
</script>
