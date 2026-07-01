<template>
	<div class="flex justify-end">
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
watch(() => props.data, (newData: any) => {
	if (newData.metric === "Throughput") {
		if (newData.data.computedThroughput === undefined) return;
		header.value = `${newData.data.computedThroughput.toFixed(4)} requests/s`
	}
}, { immediate: true })
</script>
