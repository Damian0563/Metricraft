<template>
	<NuxtLayout>
		<Navbar />
		<div class="flex-1 flex items-center justify-center">
			<div class="text-center">
				<h1 class="text-9xl font-bold text-[#00F376] tracking-widest"
					style="text-shadow: 0 0 4px #00F376, 0 0 12px #00F376;user-select: none;">
					404
				</h1>
				<p class="mt-6 text-2xl text-white">
					Page Not Found
				</p>
				<NuxtLink to="/" class="mt-10 inline-block">
					<button class="btn">Go Home</button>
				</NuxtLink>
			</div>
		</div>
	</NuxtLayout>
</template>
