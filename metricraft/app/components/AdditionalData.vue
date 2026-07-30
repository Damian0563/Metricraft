<template>
	<ClientOnly>
		<Teleport to="body">
			<AnimatePresence>
				<motion.div v-if="show" key="additional-data-overlay"
					class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm md:p-8"
					:initial="{ opacity: 0 }" :animate="{ opacity: 1 }" :exit="{ opacity: 0 }" :transition="{ duration: 0.2 }"
					@click.self="close">
					<motion.div
						class="flex w-full max-w-6xl max-h-[90vh] flex-col overflow-hidden rounded-xl bg-white shadow-xl ring-1 ring-slate-100"
						role="dialog" aria-modal="true" :initial="{ opacity: 0, scale: 0.96, y: 16 }"
						:animate="{ opacity: 1, scale: 1, y: 0 }" :exit="{ opacity: 0, scale: 0.96, y: 16 }"
						:transition="{ type: 'spring', duration: 0.35, bounce: 0.2 }" @click.stop>
						<div class="grid grid-cols-[1fr_auto_1fr] items-center gap-4 px-6 py-4">
							<button @click="createCSV"
								class="justify-self-start flex h-8 shrink-0 items-center justify-center rounded-lg px-3 text-dark-gray transition-colors hover:bg-slate-100 hover:text-[#00F376]">
								Export to CSV
							</button>
							<span class="min-w-0 truncate text-center text-dark-gray">
								{{ props.metric }}
							</span>
							<button type="button" @click="close"
								class="justify-self-end flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-dark-gray transition-colors hover:bg-slate-100 hover:text-[#00F376]"
								aria-label="Close">
								<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
									stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5">
									<line x1="18" y1="6" x2="6" y2="18"></line>
									<line x1="6" y1="6" x2="18" y2="18"></line>
								</svg>
							</button>
						</div>
						<div class="min-h-[12rem] flex-1 overflow-y-auto px-6 py-4">
							<div ref="contentRef" class="additional-data-content w-full"></div>
						</div>
					</motion.div>
				</motion.div>
			</AnimatePresence>
		</Teleport>
	</ClientOnly>
</template>

<script setup lang="ts">
import { nextTick, toRaw } from 'vue';
import { motion, AnimatePresence } from 'motion-v';

const props = defineProps<{
	show: boolean;
	data: HTMLElement | null;
	metric: string;
}>();
const emit = defineEmits<{
	close: [];
}>();
const contentRef = ref<HTMLDivElement | null>(null);
const close = () => emit('close');

const resolveSourceElement = (): HTMLElement | null => {
	if (!props.data) return null;
	const source = toRaw(props.data) as HTMLElement;
	return source.querySelector('table');
};

const renderContent = async (): Promise<void> => {
	await nextTick();
	if (!contentRef.value) return;
	contentRef.value.replaceChildren();
	const source = resolveSourceElement();
	if (!source) return;
	contentRef.value.appendChild(source.cloneNode(true));
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

watch([() => props.show, () => props.data], ([visible]) => {
	if (visible) void renderContent();
}, { flush: 'post' });

const onKeydown = (event: KeyboardEvent) => {
	if (event.key === 'Escape' && props.show) close();
};

onMounted(() => {
	if (props.show) void renderContent();
	window.addEventListener('keydown', onKeydown);
});

onUnmounted(() => {
	window.removeEventListener('keydown', onKeydown);
});
</script>

<style scoped>
.additional-data-content :deep(table) {
	display: table;
	width: 100%;
	max-width: 36rem;
	margin-inline: auto;
}
</style>
