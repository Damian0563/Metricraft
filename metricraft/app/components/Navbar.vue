<template>
	<div class="flex items-center text-white p-3 w-full h-24" style="background-color: #00F376;">
		<NuxtLink to="/" aria-label="Home">
			<img src="~/assets/src/logo.svg" class="h-36 w-36" alt="Metricraft logo" decoding="async" loading="lazy" />
		</NuxtLink>
		<div class="ml-auto gap-3">
			<NuxtLink to="/about" aria-label="About">
				<button class="btn ml-4" aria-label="About">About</button>
			</NuxtLink>
			<NuxtLink to="/manual" aria-label="Documentation">
				<button class="btn ml-4" aria-label="Documentation">Documentation</button>
			</NuxtLink>
		</div>
	</div>
</template>

<script setup lang="ts">
</script>
