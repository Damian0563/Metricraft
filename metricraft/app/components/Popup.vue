<template>


</template>

<script setup lang="ts">
const emit = defineEmits(["close"])
const props = defineProps({
	message: {
		type: String,
		required: true,
	},
})

</script>
