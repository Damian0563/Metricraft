<template>
	<div>
		<div class="flex items-center text-white p-3 w-full h-24" style="background-color: #00F376;">
			<NuxtLink to="/" aria-label="Home">
				<img src="/logo.svg" class="h-36 w-36" alt="Metricraft logo" />
			</NuxtLink>
			<div class="ml-auto flex gap-3">
				<NuxtLink to="/about" class="btn ml-4" aria-label="About">
					About
				</NuxtLink>
				<NuxtLink to="/manual" class="btn ml-4" aria-label="Documentation">
					Documentation
				</NuxtLink>
			</div>
		</div>
	</div>

</template>

<script setup lang="ts">
</script>
