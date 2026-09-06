<template>
	<ClientOnly>
		<motion.section class="customizer fixed inset-0 z-20 flex flex-col pl-20 bg-[#0a0b0a] text-white" role="dialog"
			aria-modal="true" aria-label="Customize dashboard view" :initial="{ opacity: 0, y: 12 }"
			:animate="{ opacity: 1, y: 0 }" :transition="{ type: 'spring', duration: 0.35, bounce: 0.15 }">
			<header class="flex shrink-0 items-center gap-4 border-b border-white/10 px-6 py-3">
				<div class="min-w-0">
					<h1 class="truncate text-xl font-bold" style="color: #00F376;">Customize dashboard view</h1>
					<p class="truncate text-xs text-white/45">
						Drag metrics onto the canvas, reorder them and set how wide and how tall each one sits.
					</p>
				</div>
				<span
					class="hidden shrink-0 rounded-full bg-white/5 px-3 py-1 text-xs font-medium text-white/60 ring-1 ring-white/10 sm:inline">
					{{ placed.length }} {{ placed.length === 1 ? 'graph' : 'graphs' }} placed
				</span>
				<div class="ml-auto flex shrink-0 items-center gap-2">
					<button type="button" @click="resetLayout" :disabled="!placed.length"
						class="rounded-lg px-3 py-2 text-sm font-medium text-white/70 transition-colors hover:bg-white/10 hover:text-white disabled:cursor-not-allowed disabled:opacity-35 disabled:hover:bg-transparent">
						Reset
					</button>
					<button type="button" @click="saveLayout" :disabled="!placed.length"
						class="rounded-lg bg-[#00F376] px-4 py-2 text-sm font-semibold text-gray-900 transition-colors hover:bg-[#00D96A] disabled:cursor-not-allowed disabled:opacity-35 disabled:hover:bg-[#00F376]">
						Save view
					</button>
					<button type="button" @click="close" aria-label="Close customizer"
						class="flex h-9 w-9 items-center justify-center rounded-lg text-white/60 transition-colors hover:bg-white/10 hover:text-[#00F376]">
						<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
							stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5">
							<line x1="18" y1="6" x2="6" y2="18" />
							<line x1="6" y1="6" x2="18" y2="18" />
						</svg>
					</button>
				</div>
			</header>

			<div class="flex min-h-0 flex-1">
				<main class="canvas min-h-0 flex-1 overflow-y-auto p-6" :class="{ 'canvas-active': dragging }"
					@dragover.prevent="onCanvasDragOver" @drop.prevent="onCanvasDrop">

					<div v-if="!placed.length"
						class="flex h-full min-h-[18rem] flex-col items-center justify-center gap-3 rounded-2xl border-2 border-dashed px-6 text-center transition-colors"
						:class="dragging ? 'border-[#00F376]/70 bg-[#00F376]/5' : 'border-white/12'">
						<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
							stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-10 w-10 text-white/25">
							<rect x="3" y="3" width="7" height="9" rx="1.5" />
							<rect x="14" y="3" width="7" height="5" rx="1.5" />
							<rect x="14" y="12" width="7" height="9" rx="1.5" />
							<rect x="3" y="16" width="7" height="5" rx="1.5" />
						</svg>
						<p class="text-sm font-medium text-white/70">Your canvas is empty</p>
						<p class="max-w-sm text-xs text-white/40">
							Drag a metric from the panel on the right, or click one, to drop its graph here.
						</p>
					</div>

					<div v-else class="grid grid-cols-1 items-start gap-4 md:grid-cols-2 lg:grid-cols-3">
						<template v-for="(card, index) in placed" :key="card.id">
							<div v-if="dropIndex === index" class="drop-marker"
								:class="[widthClass[draggedShape.span], heightClass[draggedShape.height]]" />
							<article draggable="true" :class="[widthClass[card.span], heightClass[card.height]]"
								class="group relative flex flex-col overflow-hidden rounded-xl bg-white text-slate-900 shadow-xl ring-1 ring-slate-100 transition-[height,opacity] duration-200"
								:style="{ opacity: dragging?.id === card.id ? 0.4 : 1 }" @dragstart="startCardDrag($event, card)"
								@dragend="endDrag" @dragover.prevent.stop="onCardDragOver($event, index)">

								<div class="flex items-center gap-2 border-b border-slate-100 px-4 py-3">
									<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"
										class="h-4 w-4 shrink-0 cursor-grab text-slate-300 active:cursor-grabbing">
										<circle cx="9" cy="6" r="1.6" />
										<circle cx="15" cy="6" r="1.6" />
										<circle cx="9" cy="12" r="1.6" />
										<circle cx="15" cy="12" r="1.6" />
										<circle cx="9" cy="18" r="1.6" />
										<circle cx="15" cy="18" r="1.6" />
									</svg>
									<h2 class="min-w-0 flex-1 truncate text-sm font-semibold">{{ card.name }}</h2>
									<span v-if="card.custom"
										class="shrink-0 rounded-full bg-[#00C263]/12 px-2 py-0.5 text-[11px] font-semibold text-[#00A854]">
										Custom
									</span>
									<span
										class="shrink-0 rounded-full bg-slate-100 px-2 py-0.5 text-[11px] font-medium tabular-nums text-slate-500">
										{{ card.timeframe }}
									</span>
									<button type="button" @click="removeCard(card.id)" :aria-label="`Remove ${card.name}`"
										class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-slate-400 transition-colors hover:bg-red-50 hover:text-red-500">
										<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
											stroke-width="2" stroke-linecap="round" class="h-4 w-4">
											<line x1="18" y1="6" x2="6" y2="18" />
											<line x1="6" y1="6" x2="18" y2="18" />
										</svg>
									</button>
								</div>

								<div class="flex min-h-0 flex-1 items-center justify-center px-4 py-3 text-[#00C263]">
									<GraphPreview :kind="card.kind" />
								</div>

								<div class="flex items-center gap-2 border-t border-slate-100 px-4 py-2">
									<div class="flex items-center gap-1" role="group" aria-label="Width">
										<span class="text-[10px] font-semibold tracking-wide text-slate-600">Width</span>
										<div class="flex items-center rounded-lg bg-slate-100 p-0.5">
											<button v-for="opt in widthOptions" :key="opt.span" type="button" @click="card.span = opt.span"
												:aria-label="`${opt.label} width`" :aria-pressed="card.span === opt.span"
												class="rounded-md px-2 py-0.5 text-[11px] font-semibold transition-colors" :class="card.span === opt.span
													? 'bg-white text-slate-800 shadow-sm'
													: 'text-slate-400 hover:text-slate-600'">
												{{ opt.label }}
											</button>
										</div>
									</div>
									<div class="flex items-center mx-3 gap-1" role="group" aria-label="Height">
										<span class="text-[10px] font-semibold tracking-wide text-slate-600">Height</span>
										<div class="flex items-center rounded-lg bg-slate-100 p-0.5">
											<button v-for="opt in heightOptions" :key="opt.height" type="button"
												@click="card.height = opt.height" :aria-label="`${opt.label} height`"
												:aria-pressed="card.height === opt.height"
												class="rounded-md px-2 py-0.5 text-[11px] font-semibold transition-colors" :class="card.height === opt.height
													? 'bg-white text-slate-800 shadow-sm'
													: 'text-slate-400 hover:text-slate-600'">
												{{ opt.label }}
											</button>
										</div>
									</div>
								</div>
							</article>
						</template>
						<div v-if="dropIndex === placed.length" class="drop-marker"
							:class="[widthClass[draggedShape.span], heightClass[draggedShape.height]]" />
					</div>
				</main>

				<aside class="flex w-72 shrink-0 flex-col border-l border-white/10 bg-[#111312] xl:w-80"
					@dragover.prevent="onPaletteDragOver" @drop.prevent="onPaletteDrop">
					<div class="shrink-0 border-b border-white/10 px-4 py-3">
						<h2 class="text-sm font-semibold text-white/85">Available metrics</h2>
						<p class="mt-0.5 text-xs text-white/40">Drag onto the canvas, or click to append.</p>
						<div class="relative mt-3">
							<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
								stroke-width="2" stroke-linecap="round"
								class="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-white/30">
								<circle cx="11" cy="11" r="7" />
								<line x1="21" y1="21" x2="16.65" y2="16.65" />
							</svg>
							<input v-model="query" type="search" placeholder="Filter metrics" aria-label="Filter metrics"
								class="w-full rounded-lg border border-white/10 bg-white/5 py-2 pr-3 pl-9 text-sm text-white placeholder:text-white/30 focus:border-[#00F376] focus:outline-none" />
						</div>
					</div>

					<div class="min-h-0 flex-1 space-y-2 overflow-y-auto p-4"
						:class="{ 'bg-red-500/5': dragging?.origin === 'canvas' }">
						<p v-if="!filtered.length" class="px-1 py-6 text-center text-xs text-white/35">
							{{ available.length ? 'No metric matches that filter.' : 'No derived metrics available yet.' }}
						</p>
						<div v-for="metric in filtered" :key="metric.name" :draggable="metric.enabled"
							@dragstart="startPaletteDrag($event, metric)" @dragend="endDrag" @click="metric.enabled && append(metric)"
							class="group flex items-center gap-3 rounded-xl border border-white/10 bg-white/5 px-3 py-3 transition-colors"
							:class="metric.enabled
								? 'cursor-grab hover:border-[#00F376]/60 hover:bg-white/10 active:cursor-grabbing'
								: 'cursor-not-allowed opacity-40'">
							<span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-[#00F376]/10 text-[#00F376]">
								<GraphPreview :kind="metric.kind" compact />
							</span>
							<span class="min-w-0 flex-1">
								<span class="flex min-w-0 items-center gap-1.5">
									<span class="truncate text-sm font-medium text-white/90">{{ metric.name }}</span>
									<span v-if="metric.custom"
										class="shrink-0 rounded-full bg-[#00F376]/15 px-1.5 py-0.5 text-[10px] font-semibold text-[#00F376]">
										Custom
									</span>
								</span>
								<span class="block text-[11px] text-white/40">
									{{ metric.enabled ? `${metric.timeframe} window` : 'Disabled in Settings' }}
								</span>
							</span>
							<span v-if="usage[metric.name]"
								class="shrink-0 rounded-full bg-[#00F376]/15 px-2 py-0.5 text-[11px] font-semibold tabular-nums text-[#00F376]">
								{{ usage[metric.name] }}
							</span>
						</div>
					</div>

					<p class="shrink-0 border-t border-white/10 px-4 py-3 text-[11px] leading-relaxed text-white/35">
						Drop a placed graph back here to remove it. Press <kbd class="rounded bg-white/10 px-1">Esc</kbd> to close.
					</p>
				</aside>
			</div>
		</motion.section>
	</ClientOnly>
</template>

<script setup lang="ts">
import { motion } from 'motion-v';
import { timeframeLabelFor } from '@/composables/helpers'
import type { CustomizableMetric, DisplayViewCard } from '@/composables/types/views';

type PreviewKind = 'map' | 'line' | 'bars' | 'donut' | 'gauge';
type PaletteEntry = { name: string; timeframe: string; enabled: boolean; custom: boolean; kind: PreviewKind };
type PlacedCard = {
	id: string; name: string; timeframe: string; kind: PreviewKind; custom: boolean;
	span: 1 | 2 | 3; height: 1 | 2 | 3;
};

const props = defineProps<{
	metrics: CustomizableMetric[];
}>();
const emit = defineEmits<{
	close: [];
	save: [value: DisplayViewCard[]];
}>();

const widthOptions = [
	{ span: 1 as const, label: 'S' },
	{ span: 2 as const, label: 'M' },
	{ span: 3 as const, label: 'L' },
];
const widthClass: Record<1 | 2 | 3, string> = {
	1: 'col-span-1',
	2: 'md:col-span-2 lg:col-span-2',
	3: 'md:col-span-2 lg:col-span-3',
};
const heightOptions = [
	{ height: 1 as const, label: 'S' },
	{ height: 2 as const, label: 'M' },
	{ height: 3 as const, label: 'L' },
];
const heightClass: Record<1 | 2 | 3, string> = {
	1: 'h-52',
	2: 'h-80',
	3: 'h-[27rem]',
};

const kindFor = (name: string): PreviewKind => {
	const key = name.toLowerCase();
	if (key.includes('geograph')) return 'map';
	if (key.includes('uptime')) return 'gauge';
	if (key.includes('distribution') || key.includes('method') || key.includes('status')) return 'donut';
	if (key.includes('latency') || key.includes('throughput') || key.includes('visitors')) return 'line';
	return 'bars';
};

const available = computed<PaletteEntry[]>(() =>
	(props.metrics ?? []).map((metric) => ({
		name: metric.name,
		timeframe: metric.timeframe || '7d',
		enabled: metric.enabled,
		custom: metric.custom,
		kind: kindFor(metric.name),
	})).sort((a, b) => Number(b.enabled) - Number(a.enabled) || a.name.localeCompare(b.name))
);

const query = ref('');
const filtered = computed(() => {
	const needle = query.value.trim().toLowerCase();
	return needle ? available.value.filter((m) => m.name.toLowerCase().includes(needle)) : available.value;
});

const placed = ref<PlacedCard[]>([]);
const usage = computed<Record<string, number>>(() =>
	placed.value.reduce<Record<string, number>>((acc, card) => {
		acc[card.name] = (acc[card.name] ?? 0) + 1;
		return acc;
	}, {})
);

let seq = 0;
const cardFrom = (metric: PaletteEntry): PlacedCard => ({
	id: `card-${++seq}`,
	name: metric.name,
	timeframe: timeframeLabelFor(metric.timeframe),
	kind: metric.kind,
	custom: metric.custom,
	span: 1,
	height: 1,
});

const dragging = ref<{ origin: 'palette' | 'canvas'; metric?: PaletteEntry; id?: string } | null>(null);
const dropIndex = ref<number | null>(null);
const draggedShape = computed<{ span: 1 | 2 | 3; height: 1 | 2 | 3 }>(() => {
	const held = dragging.value?.id
		? placed.value.find((card) => card.id === dragging.value?.id)
		: undefined;
	return { span: held?.span ?? 1, height: held?.height ?? 1 };
});

const startPaletteDrag = (event: DragEvent, metric: PaletteEntry) => {
	if (!metric.enabled) return event.preventDefault();
	dragging.value = { origin: 'palette', metric };
	dropIndex.value = placed.value.length;
	if (event.dataTransfer) {
		event.dataTransfer.effectAllowed = 'copy';
		event.dataTransfer.setData('text/plain', metric.name);
	}
};

const startCardDrag = (event: DragEvent, card: PlacedCard) => {
	dragging.value = { origin: 'canvas', id: card.id };
	if (event.dataTransfer) {
		event.dataTransfer.effectAllowed = 'move';
		event.dataTransfer.setData('text/plain', card.name);
	}
};

const endDrag = () => {
	dragging.value = null;
	dropIndex.value = null;
};

const onCardDragOver = (event: DragEvent, index: number) => {
	if (!dragging.value) return;
	if (event.dataTransfer) event.dataTransfer.dropEffect = dragging.value.origin === 'palette' ? 'copy' : 'move';
	const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
	dropIndex.value = event.clientX > rect.left + rect.width / 2 ? index + 1 : index;
};

const onCanvasDragOver = (event: DragEvent) => {
	if (!dragging.value) return;
	if (event.dataTransfer) event.dataTransfer.dropEffect = dragging.value.origin === 'palette' ? 'copy' : 'move';
	if (dropIndex.value === null) dropIndex.value = placed.value.length;
};

const onCanvasDrop = () => {
	const drag = dragging.value;
	if (!drag) return;
	const target = dropIndex.value ?? placed.value.length;
	if (drag.origin === 'palette' && drag.metric) {
		placed.value.splice(target, 0, cardFrom(drag.metric));
	} else if (drag.id) {
		const from = placed.value.findIndex((card) => card.id === drag.id);
		if (from !== -1) {
			const [card] = placed.value.splice(from, 1);
			if (card) placed.value.splice(from < target ? target - 1 : target, 0, card);
		}
	}
	endDrag();
};

const onPaletteDragOver = (event: DragEvent) => {
	if (dragging.value?.origin === 'canvas' && event.dataTransfer) event.dataTransfer.dropEffect = 'move';
	dropIndex.value = null;
};

const onPaletteDrop = () => {
	if (dragging.value?.origin === 'canvas' && dragging.value.id) removeCard(dragging.value.id);
	endDrag();
};

const append = (metric: PaletteEntry) => placed.value.push(cardFrom(metric));
const removeCard = (id: string) => {
	placed.value = placed.value.filter((card) => card.id !== id);
};
const resetLayout = () => {
	placed.value = [];
};
const saveLayout = () => emit('save', placed.value.map((card) => ({
	name: card.name,
	span: card.span,
	height: card.height,
	custom: card.custom,
})));
const close = () => emit('close');

const onKeydown = (event: KeyboardEvent) => {
	if (event.key === 'Escape') close();
};
onMounted(() => {
	window.addEventListener('keydown', onKeydown);
	document.body.style.overflow = 'hidden';
});
onBeforeUnmount(() => {
	window.removeEventListener('keydown', onKeydown);
	document.body.style.overflow = '';
});
</script>

<style scoped>
.canvas {
	background-color: #0d0f0e;
	background-image: radial-gradient(circle, rgba(255, 255, 255, 0.07) 1px, transparent 1px);
	background-size: 22px 22px;
	transition: box-shadow 0.2s ease;
}

.canvas-active {
	box-shadow: inset 0 0 0 2px rgba(0, 243, 118, 0.35);
}

.drop-marker {
	border: 2px dashed rgba(0, 243, 118, 0.7);
	border-radius: 0.75rem;
	background-color: rgba(0, 243, 118, 0.07);
}

kbd {
	font-family: inherit;
	font-size: 10px;
}

@media (prefers-reduced-motion: reduce) {
	.canvas {
		transition: none;
	}
}
</style>
