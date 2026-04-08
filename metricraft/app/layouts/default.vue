<template>
	<div class="layout bg-black" :class="overflow ? 'h-auto' : 'h-screen'">
		<slot />
	</div>
</template>


<script setup lang="ts">

const overflow = ref(false)
const checkOverflow = () => {
	overflow.value = document.documentElement.scrollHeight > window.innerHeight
}
watch(overflow, () => {
	checkOverflow()
})
onMounted(() => {
	checkOverflow()
	window.addEventListener('resize', checkOverflow)
})
onUnmounted(() => {
	window.removeEventListener('resize', checkOverflow)
})
</script>
