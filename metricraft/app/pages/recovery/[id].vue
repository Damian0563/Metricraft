<template>
	<div>
		<ClientOnly>
			<Spinner :loading="loading" />
		</ClientOnly>
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<Notice :message="noticeMessage" title="Secret key updated" @close="handleNoticeClose" />
		<div class="flex min-h-[calc(100vh-96px)] items-center justify-center px-4 py-10">
			<div class="w-full max-w-xl rounded-xl bg-white p-8 shadow-lg">
				<div class="mb-6 flex justify-center">
					<div class="flex h-20 w-20 items-center justify-center rounded-2xl bg-[#00F376]/15">
						<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#00F376" stroke-width="2"
							stroke-linecap="round" stroke-linejoin="round" class="h-10 w-10">
							<path d="M12 17h.01" />
							<path d="M7 11V7a5 5 0 0 1 10 0v4" />
							<rect x="5" y="11" width="14" height="10" rx="2" />
						</svg>
					</div>
				</div>
				<div class="mb-8 text-center">
					<p class="mb-2 text-sm font-semibold uppercase tracking-wider text-[#00A652]">Account recovery</p>
					<h1 class="text-3xl font-bold text-gray-900">Create a new secret key</h1>
					<p class="mt-3 text-sm text-gray-500">
						Choose a strong secret key for your Metricraft account. Recovery links expire after 10 minutes.
					</p>
				</div>
				<div v-if="!recoveryId" class="rounded-lg border border-red-100 bg-red-50 p-4 text-sm text-red-600">
					This recovery link is missing its recovery token. Please request a new link from the sign in page.
				</div>
				<form v-else @submit.prevent="handleRecovery" class="flex flex-col gap-5">
					<div>
						<label for="secret" class="block text-sm font-medium text-gray-700 mb-1">New Secret Key</label>
						<div class="flex items-center gap-2">
							<input id="secret" :type="showSecret ? 'text' : 'password'" v-model="secret"
								placeholder="Enter your new secret key" autocomplete="new-password"
								class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition" />
							<button type="button" @click="showSecret = !showSecret"
								class="shrink-0 text-[#00F376] hover:text-[#00D96A] transition-colors"
								aria-label="Toggle secret visibility">
								<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
									stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5">
									<path v-if="!showSecret"
										d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
									<line v-if="!showSecret" x1="1" y1="1" x2="23" y2="23" />
									<template v-if="showSecret">
										<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
										<circle cx="12" cy="12" r="3" />
									</template>
								</svg>
							</button>
						</div>
					</div>
					<div>
						<label for="confirmSecret" class="block text-sm font-medium text-gray-700 mb-1">Confirm Secret Key</label>
						<div class="flex items-center gap-2">
							<input id="confirmSecret" :type="showConfirmSecret ? 'text' : 'password'" v-model="confirmSecret"
								placeholder="Confirm your new secret key" autocomplete="new-password"
								class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition" />
							<button type="button" @click="showConfirmSecret = !showConfirmSecret"
								class="shrink-0 text-[#00F376] hover:text-[#00D96A] transition-colors"
								aria-label="Toggle confirmation visibility">
								<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
									stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5">
									<path v-if="!showConfirmSecret"
										d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
									<line v-if="!showConfirmSecret" x1="1" y1="1" x2="23" y2="23" />
									<template v-if="showConfirmSecret">
										<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
										<circle cx="12" cy="12" r="3" />
									</template>
								</svg>
							</button>
						</div>
						<p v-if="confirmSecret.length > 0" class="mt-2 text-sm font-medium"
							:class="secretsMatch ? 'text-green-600' : 'text-red-500'">
							{{ secretsMatch ? 'Secret keys match.' : 'Secret keys do not match.' }}
						</p>
					</div>
					<div class="min-h-12">
						<div v-if="(secret && confirmSecret.length === 0) || secretsMatch" class="flex flex-col gap-2">
							<span class="text-sm font-medium"
								:class="passwordStrength ? (passwordStrength <= 1 ? 'text-red-500' : passwordStrength <= 2 ? 'text-orange-500' : passwordStrength <= 3 ? 'text-green-700' : 'text-green-500') : 'text-gray-200'">
								{{ passwordStrengthMessages[passwordStrength - 1] }}
							</span>
							<div class="flex h-2 gap-1">
								<div v-for="i in 4" :key="i" class="flex-1 rounded-full transition-all duration-200"
									:class="i <= passwordStrength ? (i <= 1 ? 'bg-red-500' : i <= 2 ? 'bg-orange-500' : i <= 3 ? 'bg-green-700' : 'bg-green-500') : 'bg-gray-200'">
								</div>
							</div>
						</div>
					</div>
					<button type="submit" :disabled="loading || !canSubmit"
						class="w-full cursor-pointer py-3 mt-2 text-black font-semibold bg-[#00F376] hover:text-white rounded-lg shadow-lg hover:bg-black transition delay-100 ease-in-out disabled:opacity-60 disabled:cursor-not-allowed disabled:hover:bg-[#00F376] disabled:hover:text-black">
						{{ loading ? 'Updating...' : 'Update Secret Key' }}
					</button>
				</form>
				<button type="button" @click="navigateTo('/')"
					class="mx-auto mt-6 block text-sm font-medium text-gray-500 hover:text-[#00F376] transition-colors">
					Back to sign in
				</button>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { checkRecovery } from '~/calls/welcome';
import { evaluatePasswordStrength } from '@/composables/helpers';
definePageMeta({
	layout: 'entry',
})
const route = useRoute()
const recoveryId = computed(() => {
	const id = route.params.id
	return Array.isArray(id) ? id[0] : id
})
const secret = ref('')
const confirmSecret = ref('')
const showSecret = ref(false)
const showConfirmSecret = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const noticeMessage = ref('')
const passwordStrength = computed(() => evaluatePasswordStrength(secret.value) || 1)
const secretsMatch = computed(() => secret.value === confirmSecret.value && secret.value.length > 0)
const canSubmit = computed(() => Boolean(secret.value && confirmSecret.value && secretsMatch.value && passwordStrength.value > 1))
const passwordStrengthMessages = [
	'Weak password, try something more complex',
	'This password\'s stength is good enough.',
	'This password\'s strength is okay.',
	'This password\'s strength is excellent.'
]

const handleNoticeClose = () => {
	noticeMessage.value = ''
	setTimeout(() => {
		navigateTo('/')
	}, 1500)
}

const handleRecovery = async () => {
	if (!recoveryId.value) {
		errorMessage.value = 'This recovery link is invalid.'
		return
	}
	if (!secret.value || !confirmSecret.value) {
		errorMessage.value = 'Please fill in all fields.'
		return
	}
	if (!secretsMatch.value) {
		errorMessage.value = 'Secret keys do not match.'
		return
	}
	if (passwordStrength.value === 1) {
		errorMessage.value = 'Password is too weak.'
		return
	}
	loading.value = true
	try {
		await checkRecovery(recoveryId.value, secret.value)
		noticeMessage.value = 'Your secret key has been updated. You can now sign in with the new key.'
		secret.value = ''
		confirmSecret.value = ''
	} catch (e) {
		errorMessage.value = String(e)
	} finally {
		loading.value = false
	}
}
</script>
