<template>
	<Transition name="scrim">
		<div v-if="options" class="nav-scrim fixed inset-0 z-30 bg-black/30" aria-hidden="true" @click="options = false" />
	</Transition>
	<aside
		class="nav-sidebar fixed top-0 left-0 z-40 flex flex-col text-white text-2xl font-bold h-screen py-4 px-3 duration-300 ease-in-out"
		:class="options ? 'w-64' : 'w-20'" style="background-color: #00F376;">
		<nav class="flex flex-col h-full gap-8">
			<div class="flex items-center" :class="options ? 'justify-between' : 'justify-center'">
				<div class="flex items-center overflow-hidden" :class="{ 'nav-collapsed': !options }">
					<img src="/logo.svg" class="h-8 nav-label cursor-pointer shrink-0" alt="logo, go back to main view"
						@click="getBack" />
				</div>
				<button
					class="nav-menu-btn bg-white text-black rounded-full p-2 text-sm hover:cursor-pointer transition-shadow duration-300 shrink-0"
					aria-label="Toggle navigation" @click="options = !options">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
					</svg>
				</button>
			</div>
			<div class="flex flex-col gap-1">
				<button type="button" class="nav-item" :class="{ 'nav-collapsed': !options }" title="Home" @click="getBack"
					aria-label="Home">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0" fill="none" viewBox="0 0 24 24"
						stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
							d="M3 12l9-9 9 9M5 10v10a1 1 0 001 1h3v-6h6v6h3a1 1 0 001-1V10" />
					</svg>
					<span class="nav-label">Home</span>
				</button>
				<button type="button" class="nav-item" :class="{ 'nav-collapsed': !options }" title="Team"
					@click="navigateTo('/invite'); options = false">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0" fill="none" viewBox="0 0 24 24"
						stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
							d="M17 20h5v-2a4 4 0 00-3-3.87M9 20H4v-2a4 4 0 013-3.87m6-1.13a4 4 0 10-4-4 4 4 0 004 4z" />
					</svg>
					<span class="nav-label">Team</span>
				</button>
				<button type="button" class="nav-item" :class="{ 'nav-collapsed': !options }" title="Metricraft workers"
					@click="navigateTo('/workers'); options = false">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0" fill="none" viewBox="0 0 24 24"
						stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
							d="M11.42 15.17L17.25 21A2.652 2.652 0 0021 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 11-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 004.486-6.336l-3.276 3.277a3.004 3.004 0 01-2.25-2.25l3.276-3.276a4.5 4.5 0 00-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437l1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008z" />
					</svg>
					<span class="nav-label">Metricraft workers</span>
				</button>
				<button type="button" class="nav-item" :class="{ 'nav-collapsed': !options }" title="Settings"
					@click="navigateTo('/dashboard?settings'); options = false">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0" fill="none" viewBox="0 0 24 24"
						stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
							d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
							d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
					</svg>
					<span class="nav-label">Settings</span>
				</button>
			</div>

			<button
				class="sign-out-btn mt-auto text-sm font-medium text-white bg-red-500 rounded-full px-3 py-2 hover:bg-red-600"
				:class="{ 'nav-collapsed': !options }" title="Sign out" @click="signOut">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0" fill="none" viewBox="0 0 24 24"
					stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
						d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
				</svg>
				<span class="nav-label">Sign out</span>
			</button>
		</nav>
	</aside>
</template>

<script setup lang="ts">
import { invalidateCookie } from "~/composables/helpers"
const options = ref(false)
const signOut = () => {
	invalidateCookie()
	navigateTo('/')
}
const router = useRouter()
const getBack = () => {
	options.value = false
	if (router.currentRoute.value.name === 'dashboard') {
		router.replace({ query: {} })
	} else {
		navigateTo('/dashboard')
	}
}
const onKeydown = (e: KeyboardEvent) => {
	if (e.key === 'Escape') options.value = false
}
onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))
</script>

<style scoped>
.nav-sidebar {
	box-shadow:
		2px 0 6px rgba(0, 0, 0, 0.06),
		4px 0 18px rgba(0, 0, 0, 0.08);
	/* Only animate width, and hint the compositor so it isn't recalculated blindly */
	transition-property: width;
	will-change: width;
	/* Prevent the sidebar's internal layout from invalidating on every frame */
	contain: layout paint;
}

.nav-item,
.sign-out-btn {
	display: flex;
	align-items: center;
	justify-content: flex-start;
	width: 100%;
}

.nav-item {
	padding: 0.625rem 0.75rem;
	font-size: 0.95rem;
	font-weight: 600;
	color: rgb(17 24 39);
	border-radius: 0.75rem;
	cursor: pointer;
	background: transparent;
	border: none;
	transition: background-color 0.2s ease, color 0.2s ease;
}

.nav-item:hover {
	background-color: rgba(255, 255, 255, 0.55);
}

.nav-item:active {
	background-color: rgba(255, 255, 255, 0.75);
}

/* Labels stay mounted and fade/collapse instead of using v-if (avoids reflow pops) */
.nav-label {
	white-space: nowrap;
	overflow: hidden;
	opacity: 1;
	max-width: 12rem;
	margin-left: 0.75rem;
	transition: opacity 0.2s ease, max-width 0.3s ease, margin-left 0.3s ease;
}

.nav-collapsed {
	justify-content: center;
}

.nav-collapsed .nav-label {
	opacity: 0;
	max-width: 0;
	margin-left: 0;
}

.scrim-enter-active,
.scrim-leave-active {
	transition: opacity 0.3s ease;
}

.scrim-enter-from,
.scrim-leave-to {
	opacity: 0;
}

@media (prefers-reduced-motion: reduce) {

	.nav-sidebar,
	.nav-label {
		transition: none;
	}
}

.nav-menu-btn {
	box-shadow:
		0 1px 2px rgba(0, 0, 0, 0.08),
		0 4px 14px rgba(0, 0, 0, 0.12);
}

.nav-menu-btn:hover {
	box-shadow:
		0 2px 4px rgba(0, 0, 0, 0.1),
		0 8px 24px rgba(0, 0, 0, 0.16);
}

.nav-menu-btn:active {
	box-shadow:
		0 1px 2px rgba(0, 0, 0, 0.1),
		0 2px 6px rgba(0, 0, 0, 0.12);
}

.sign-out-btn {
	box-shadow:
		0 2px 6px rgba(239, 68, 68, 0.35),
		0 1px 3px rgba(185, 28, 28, 0.2);
}

.sign-out-btn:hover {
	box-shadow:
		0 4px 14px rgba(239, 68, 68, 0.45),
		0 2px 6px rgba(185, 28, 28, 0.25);
}

.sign-out-btn:active {
	box-shadow:
		0 1px 3px rgba(239, 68, 68, 0.3),
		0 1px 2px rgba(185, 28, 28, 0.2);
}
</style>
