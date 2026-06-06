<template>
	<div>
		<DashboardNav :appName="appName" @settings="handleSettings" />
		<Popup :message="errorMessage" @close="errorMessage = ''" />
		<div class="w-full px-8 py-2">
			<div class="relative flex items-center justify-center mb-6">
				<button @click="goBack"
					class="absolute left-0 text-white hover:text-[#00F376] transition-colors duration-200 flex items-center gap-2 cursor-pointer">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
					</svg>
					Back to Dashboard
				</button>
				<h1 class="text-3xl font-bold text-center" style="color: #00F376;">Invite Team Members</h1>
			</div>
			<div class="max-w-4xl mx-auto">
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<div class="flex gap-4 mb-8 border-b border-gray-200">
						<button @click="mode = 'manual'"
							:class="['pb-4 px-2 font-medium transition-all duration-200 cursor-pointer', mode === 'manual' ? 'border-b-2 border-[#00F376] text-[#00F376]' : 'text-gray-500 hover:text-gray-700']">
							Manual Entry
						</button>
						<button @click="mode = 'batch'"
							:class="['pb-4 px-2 font-medium transition-all duration-200 cursor-pointer', mode === 'batch' ? 'border-b-2 border-[#00F376] text-[#00F376]' : 'text-gray-500 hover:text-gray-700']">
							Batch Upload (CSV)
						</button>
					</div>
					<div v-if="mode === 'manual'" class="space-y-6">
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">Email Address</label>
							<div class="flex gap-3">
								<input v-model="emailInput" type="email" placeholder="colleague@example.com"
									class="flex-1 px-4 py-2 rounded-lg border border-gray-200 focus:outline-none focus:border-[#00F376] transition-colors text-gray-800"
									@keyup.enter="addEmail" />
								<button @click="addEmail"
									class="px-6 py-2 bg-[#00F376] text-gray-900 font-semibold rounded-lg hover:bg-[#00D96A] transition-colors shadow-md cursor-pointer">
									Add
								</button>
							</div>
						</div>
						<div v-if="emails.length > 0">
							<h3 class="text-sm font-semibold text-gray-800 mb-3">Pending Invites ({{ emails.length }})</h3>
							<div class="space-y-2 max-h-60 overflow-y-auto pr-2">
								<div v-for="(email, index) in emails" :key="index"
									class="flex items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-100">
									<span class="text-gray-700">{{ email }}</span>
									<button @click="removeEmail(index)" class="text-red-500 hover:text-red-700 cursor-pointer">
										<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
											<path fill-rule="evenodd"
												d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
												clip-rule="evenodd" />
										</svg>
									</button>
								</div>
							</div>
						</div>
					</div>
					<div v-if="mode === 'batch'" class="space-y-6">
						<div
							class="border-2 border-dashed border-gray-200 rounded-xl p-12 text-center hover:border-[#00F376] transition-colors cursor-pointer group"
							@click="fileInput?.click()">
							<input type="file" ref="fileInput" class="hidden" accept=".csv" @change="handleFileUpload" />
							<div class="flex flex-col items-center">
								<svg xmlns="http://www.w3.org/2000/svg"
									class="h-12 w-12 text-gray-400 mb-4 group-hover:text-[#00F376] transition-colors" fill="none"
									viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
										d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
								</svg>
								<p class="text-lg font-medium text-gray-700">Click to upload CSV</p>
								<p class="text-sm text-gray-500 mt-1">File should contain a column with email addresses</p>
							</div>
						</div>
						<div v-if="csvFile"
							class="flex items-center justify-between p-4 bg-[#00F376]/10 rounded-lg border border-[#00F376]/20">
							<div class="flex items-center gap-3">
								<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-[#00F376]" fill="none" viewBox="0 0 24 24"
									stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
										d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
								</svg>
								<div class="text-left">
									<p class="text-sm font-medium text-gray-800">{{ csvFile.name }}</p>
									<p class="text-xs text-gray-500">{{ (csvFile.size / 1024).toFixed(2) }} KB</p>
								</div>
							</div>
							<button @click="csvFile = null" class="text-gray-400 hover:text-gray-600 cursor-pointer">
								<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
									<path fill-rule="evenodd"
										d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
										clip-rule="evenodd" />
								</svg>
							</button>
						</div>
					</div>
					<div class="mt-10 pt-6 border-t border-gray-100 flex justify-end">
						<button @click="sendInvites" :disabled="!canSend"
							class="px-8 py-3 bg-[#00F376] text-gray-900 font-bold rounded-lg hover:bg-[#00D96A] transition-all shadow-lg disabled:opacity-50 disabled:cursor-not-allowed uppercase tracking-wider text-sm cursor-pointer">
							Send Invites
						</button>
					</div>
				</div>
			</div>
			<div class="max-w-6xl mx-auto mt-8 grid gap-8 lg:grid-cols-2">
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<div class="mb-6">
						<h2 class="text-2xl font-bold text-gray-900">Current Team</h2>
						<p class="text-sm text-gray-500 mt-1">Users who currently have access to this project.</p>
					</div>
					<div class="space-y-3">
						<div v-for="user in teamUsers" :key="user.mail"
							class="flex flex-col gap-3 p-4 bg-gray-50 rounded-lg border border-gray-100 xl:flex-row xl:items-center xl:justify-between">
							<div class="flex min-w-0 items-center gap-3">
								<div class="h-10 w-10 shrink-0 rounded-full bg-[#00F376]/10 flex items-center justify-center text-[#00F376] font-semibold">
									{{ user.initials }}
								</div>
								<div class="min-w-0">
									<p class="font-medium text-gray-900 break-all">{{ user.mail }}</p>
									<p class="text-xs text-gray-500">{{ user.role }}</p>
								</div>
							</div>
							<span class="w-fit shrink-0 px-3 py-1 rounded-full bg-[#00F376]/10 text-[#00A652] text-xs font-semibold uppercase tracking-wide">
								Active
							</span>
						</div>
					</div>
				</div>
				<div class="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
					<div class="mb-6">
						<h2 class="text-2xl font-bold text-gray-900">Pending Verification</h2>
						<p class="text-sm text-gray-500 mt-1">Review users waiting for access to this project.</p>
					</div>
					<div v-if="pendingUserList.length > 0" class="space-y-3">
						<div v-for="user in pendingUserList" :key="user.mail"
							class="flex flex-col gap-3 p-4 bg-gray-50 rounded-lg border border-gray-100 xl:flex-row xl:items-center xl:justify-between">
							<div class="flex min-w-0 items-center gap-3">
								<div class="h-10 w-10 shrink-0 rounded-full bg-gray-200 flex items-center justify-center text-gray-500">
									<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 20 20" fill="currentColor">
										<path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
											clip-rule="evenodd" />
									</svg>
								</div>
								<div class="min-w-0">
									<p class="font-medium text-gray-900 break-all">{{ user.mail }}</p>
									<p class="text-xs text-gray-500">Awaiting verification</p>
								</div>
							</div>
							<div class="flex shrink-0 items-center gap-2">
								<button @click="handlePendingUser(user.mail, true)"
									class="px-4 py-2 bg-[#00F376] text-gray-900 font-semibold rounded-lg hover:bg-[#00D96A] transition-colors cursor-pointer">
									Accept
								</button>
								<button @click="handlePendingUser(user.mail, false)"
									class="px-4 py-2 bg-red-50 text-red-600 font-semibold rounded-lg hover:bg-red-100 transition-colors cursor-pointer">
									Reject
								</button>
							</div>
						</div>
					</div>
					<p v-else class="text-sm text-gray-500">No users are pending verification.</p>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { getCookie, validateEmail } from "@/composables/helpers";
import { getPendingUsers, handlePermissionDecision } from "@/calls/invite";
const appName = useState<string>('appName');
const mode = ref<'manual' | 'batch'>('manual')
const emailInput = ref('')
const emails = ref<string[]>([])
const csvFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const errorMessage = ref('')
const teamUsers = ref([
	{ mail: 'alex@metricraft.dev', role: 'Owner', initials: 'AL' },
	{ mail: 'sam@metricraft.dev', role: 'Developer', initials: 'SA' },
	{ mail: 'jordan@metricraft.dev', role: 'Analyst', initials: 'JO' },
])
const { data: pendingUsers, error: pendingUsersError } = await useAsyncData('pendingUsers', () => getPendingUsers(), { default: () => [] })
const pendingUserList = computed(() => pendingUsers.value ?? [])

const addEmail = () => {
	const email = emailInput.value.trim()
	if (emails.value.includes(email)) {
		errorMessage.value = 'This email address is already in the list of invites.'
	}
	else if (email && validateEmail(email)) {
		emails.value.push(email)
		emailInput.value = ''
	} else {
		errorMessage.value = 'Please enter a valid email address.'
		emailInput.value = ''
	}
}

const removeEmail = (index: number) => {
	emails.value.splice(index, 1)
}

const handleFileUpload = (event: Event) => {
	const target = event.target as HTMLInputElement
	if (target.files && target.files[0]) {
		csvFile.value = target.files[0]
	}
}

const canSend = computed(() => {
	if (mode.value === 'manual') return emails.value.length > 0
	return csvFile.value !== null
})

const sendInvites = () => {
	if (mode.value === 'manual') emails.value = []
	else csvFile.value = null
}

const removePendingUser = (mail: string) => {
	pendingUsers.value = pendingUserList.value.filter(user => user.mail !== mail)
}

const handlePendingUser = (mail: string, action: boolean) => {
	removePendingUser(mail)
	handlePermissionDecision(mail, action)
}

const handleSettings = () => {
	navigateTo('/dashboard?settings')
}

const goBack = () => {
	navigateTo('/dashboard')
}

onMounted(() => {
	const cookie = getCookie("session-token");
	if (!cookie) {
		navigateTo("/");
	}
})
</script>
