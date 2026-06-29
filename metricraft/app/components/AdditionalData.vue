<template>
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm md:p-8 w-full h-full"
		@click.self="close">
		<div
			class="relative flex w-full max-w-xl max-h-[90vh] flex-col overflow-hidden rounded-xl bg-white shadow-xl ring-1 ring-slate-100"
			role="dialog" aria-modal="true">
			<button @click="createCSV"
				class="absolute top-3 left-3 z-10 flex h-8 w-32 items-center justify-center rounded-lg text-dark-gray transition-colors hover:bg-slate-100 hover:text-[#00F376]">
				Export to CSV </button>
			<span class="absolute top-3 left-38 z-10 flex h-8 w-72 items-center justify-center rounded-lg text-dark-gray">
				{{ props.metric }}
			</span>
			<button type="button" @click="close"
				class="absolute top-3 right-3 z-10 flex h-8 w-8 items-center justify-center rounded-lg text-dark-gray transition-colors hover:bg-slate-100 hover:text-[#00F376]"
				aria-label="Close">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
					stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5">
					<line x1="18" y1="6" x2="6" y2="18"></line>
					<line x1="6" y1="6" x2="18" y2="18"></line>
				</svg>
			</button>
			<div class="overflow-y-auto px-6 py-8 pt-10 mt-4">
				<div ref="contentRef" class="additional-data-content"></div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
const props = defineProps<{
	data: HTMLElement | null;
	metric: string;
}>();
const emit = defineEmits<{
	close: [];
}>();
const contentRef = ref<HTMLDivElement | null>(null);
const close = () => emit('close');
const renderContent = () => {
	if (!contentRef.value) return;
	contentRef.value.replaceChildren();
	if (!props.data) return;
	for (const child of Array.from(props.data.childNodes)) {
		contentRef.value.appendChild(child.cloneNode(true));
	}
};

const escapeCsvField = (value: string): string => {
	const normalized = value.replace(/\r?\n/g, ' ').trim();
	if (/[",\n\r]/.test(normalized)) {
		return `"${normalized.replace(/"/g, '""')}"`;
	}
	return normalized;
};

const createCSV = () => {
	if (!contentRef.value) return;
	const table = contentRef.value.querySelector('table');
	if (!table) return;
	const rows = Array.from(table.querySelectorAll('tr')).map((row) =>
		Array.from(row.querySelectorAll('th, td')).map((cell) =>
			escapeCsvField(cell.textContent ?? '')
		)
	);
	if (!rows.length) return;
	const csv = rows.map((row) => row.join(',')).join('\n');
	const blob = new Blob([csv], { type: 'text/csv;charset=utf-8' });
	const url = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = url;
	link.download = `${props.metric}.csv`;
	document.body.appendChild(link);
	link.click();
	document.body.removeChild(link);
	URL.revokeObjectURL(url);
};

onMounted(() => {
	renderContent();
	window.addEventListener('keydown', onKeydown);
});

onUnmounted(() => {
	window.removeEventListener('keydown', onKeydown);
});

watch(() => props.data, renderContent);

const onKeydown = (event: KeyboardEvent) => {
	if (event.key === 'Escape') close();
};
</script>

<style scoped>
.additional-data-content {
	display: flex;
	justify-content: center;
	width: 100%;
}

.additional-data-content :deep(table) {
	width: 100%;
	max-width: 36rem;
}
</style>
