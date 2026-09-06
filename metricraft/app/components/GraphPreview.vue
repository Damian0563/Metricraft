<template>
	<svg xmlns="http://www.w3.org/2000/svg" :viewBox="compact ? '0 0 40 40' : '0 0 200 96'"
		:class="compact ? 'h-5 w-5' : 'h-full max-h-56 w-full'" fill="none" aria-hidden="true" role="presentation">
		<g :transform="compact ? 'translate(4 12.32) scale(0.16)' : undefined" stroke="currentColor"
			:stroke-width="compact ? 14 : 2.5" stroke-linecap="round" stroke-linejoin="round">

			<template v-if="kind === 'bars'">
				<rect v-for="(h, i) in [38, 62, 30, 78, 52, 88, 44]" :key="i" :x="12 + i * 26" :y="90 - h" width="14"
					:height="h" rx="3" fill="currentColor" stroke="none" :opacity="0.25 + (i % 3) * 0.28" />
			</template>

			<template v-else-if="kind === 'line'">
				<path d="M8 74 L40 58 L72 66 L104 34 L136 46 L168 20 L192 30 L192 90 L8 90 Z" fill="currentColor" opacity="0.12"
					stroke="none" />
				<polyline points="8,74 40,58 72,66 104,34 136,46 168,20 192,30" />
				<circle v-for="p in [[104, 34], [168, 20]]" :key="p.join()" :cx="p[0]" :cy="p[1]" r="4" fill="currentColor"
					stroke="none" />
			</template>

			<template v-else-if="kind === 'donut'">
				<ellipse cx="100" cy="48" rx="54" ry="34" opacity="0.15" stroke-width="16" />
				<ellipse cx="100" cy="48" rx="54" ry="34" stroke-width="16" stroke-dasharray="112 168" />
				<ellipse cx="100" cy="48" rx="54" ry="34" stroke-width="16" stroke-dasharray="62 218"
					stroke-dashoffset="-112" opacity="0.5" />
			</template>

			<template v-else-if="kind === 'gauge'">
				<path d="M52 76 A48 48 0 1 1 148 76" opacity="0.18" stroke-width="14" />
				<path d="M52 76 A48 48 0 0 1 128 34" stroke-width="14" />
			</template>

			<template v-else>
				<ellipse cx="100" cy="48" rx="84" ry="40" opacity="0.35" />
				<ellipse cx="100" cy="48" rx="56" ry="40" opacity="0.18" />
				<ellipse cx="100" cy="48" rx="28" ry="40" opacity="0.35" />
				<line x1="16" y1="48" x2="184" y2="48" opacity="0.35" />
				<circle cx="66" cy="30" r="5" fill="currentColor" stroke="none" />
				<circle cx="128" cy="62" r="5" fill="currentColor" stroke="none" opacity="0.6" />
				<circle cx="104" cy="22" r="3.5" fill="currentColor" stroke="none" opacity="0.45" />
			</template>
		</g>
	</svg>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
	kind: 'map' | 'line' | 'bars' | 'donut' | 'gauge';
	compact?: boolean;
}>(), { compact: false });
</script>
