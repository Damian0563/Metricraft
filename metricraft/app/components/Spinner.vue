<template>
	<div>
		<div class="spinner" v-if="props.loading">
			<div class="loader"></div>
		</div>
	</div>
</template>

<script setup lang="ts">
const props = defineProps<{
	loading: boolean
}>()
</script>

<style scoped>
.spinner {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 100;
	background: rgba(0, 0, 0, 0.5);
}

.loader {
	width: 50px;
	height: 50px;
	border: 4px solid rgba(0, 243, 118, 0.2);
	border-top-color: #00f376;
	border-radius: 50%;
	animation: spin 1s linear infinite;
}

@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}
</style>
