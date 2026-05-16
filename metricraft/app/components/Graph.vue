<template>
	<div class="flex flex-col rounded shadow-lg bg-white w-full h-48 md:h-64 lg:h-96 text-black aspect-square">
		<h1 class="text-2xl font-bold text-center">{{ props.name }}</h1>
	</div>
</template>

<script setup lang="ts">
const props = defineProps<{
	name: string;
	data: any;
}>();

console.log(props.data)

</script>
