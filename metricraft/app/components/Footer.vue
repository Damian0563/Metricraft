<template>
	<footer class="bg-gray-900 text-white py-6 px-4 mt-auto w-full">
		<div class="max-w-7xl mx-auto flex flex-col md:flex-row justify-between items-center gap-4">
			<div class="text-sm text-gray-400">
				© {{ new Date().getFullYear() }} Metricraft. All rights reserved.
			</div>
			<div class="flex gap-6">
				<NuxtLink to="/about" class="text-sm text-gray-400 hover:text-white transition-colors">
					About
				</NuxtLink>
				<NuxtLink to="/manual" class="text-sm text-gray-400 hover:text-white transition-colors">
					Documentation
				</NuxtLink>
				<NuxtLink to="/privacy" class="text-sm text-gray-400 hover:text-white transition-colors">
					Privacy
				</NuxtLink>
			</div>
		</div>
	</footer>
</template>

<script setup lang="ts">
</script>

<style scoped>
footer {
	position: fixed;
	bottom: 0;
	left: 0;
}
</style>