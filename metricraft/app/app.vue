<template>
	<NuxtLayout>
		<NuxtPage />
	</NuxtLayout>
</template>

<style>
.page-enter-active,
.page-leave-active {
	transition:
		opacity 0.38s cubic-bezier(0.22, 1, 0.36, 1),
		transform 0.38s cubic-bezier(0.22, 1, 0.36, 1),
		filter 0.38s cubic-bezier(0.22, 1, 0.36, 1);
}

.page-enter-from {
	opacity: 0;
	transform: translateY(14px) scale(0.988);
	filter: blur(5px);
}

.page-leave-to {
	opacity: 0;
	transform: translateY(-10px) scale(0.992);
	filter: blur(3px);
}

@media (prefers-reduced-motion: reduce) {

	.page-enter-active,
	.page-leave-active {
		transition: opacity 0.15s ease;
	}

	.page-enter-from,
	.page-leave-to {
		transform: none;
		filter: none;
	}
}
</style>
