<template>
	<div class="flex justify-center items-center py-16 px-4">
		<div class="bg-white p-8 shadow-lg w-full max-w-lg rounded-xl">
			<div class="flex justify-center mb-6">
				<img src="/favicon.ico" alt="logo" class="w-24 h-24 rounded-lg" decoding="async" loading="lazy" />
			</div>
			<h2 class="text-2xl font-bold text-center text-gray-800 mb-6">
				{{ oldUser ? 'Welcome Back' : 'Create Account' }}
			</h2>
			<form @submit.prevent class="flex flex-col gap-5">
				<div v-if="!oldUser">
					<label for="appName" class="block text-sm font-medium text-gray-700 mb-1">App Name</label>
					<input id="appName" type="text" v-model="appName" placeholder="Enter your app name"
						class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition" />
				</div>
				<div class="relative">
					<div class="flex items-center gap-2">
						<label for="email" class="block text-sm font-medium text-gray-700">Email</label>
						<div class="group relative">
							<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#00F376" class="w-5 h-5 text-gray-400">
								<path fill-rule="evenodd"
									d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12Zm8.706-1.442c1.146-.573 2.437.463 2.126 1.706l-.709 2.836.042-.02a.75.75 0 01.67 1.34l-.04.022c-1.147.573-2.438-.463-2.127-1.706l.71-2.836-.042.02a.75.75 0 11-.671-1.34l.041-.022ZM12 9a.75.75 0 100-1.5.75.75 0 000 1.5Z"
									clip-rule="evenodd" />
							</svg>
							<div
								class="absolute right-0 bottom-full mb-2 hidden group-hover:block w-96 p-3 bg-[#00F376] text-black text-s rounded-lg shadow-lg z-10">
								<div v-if="!oldUser">Your email is used to send you a recovery link if you lose your secret key. We
									don't
									use it for anything else. For more information refer to our documentation.</div>
								<div v-else>Enter the email associated with your account to sign in.</div>
								<div class="absolute right-2 -bottom-1 w-2 h-2 bg-gray-900 rotate-45"></div>
							</div>
						</div>
					</div>
					<input id="email" type="email" v-model="mail" placeholder="you@example.com"
						class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition mt-1" />
				</div>
				<div>
					<label for="secret" class="block text-sm font-medium text-gray-700 mb-1">Secret Key</label>
					<div class="w-full flex">
						<input id="secret" :type="showSecret ? 'text' : 'password'" v-model="secret"
							placeholder="Enter your secret key" :keyup="oldUser ? handleTextInput : null"
							class="w-96 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition" />
						<svg @click="showSecret = !showSecret" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
							stroke="#00F376" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
							class="w-5 h-5 text-gray-400 my-auto ml-2 cursor-pointer">
							<path v-if="!showSecret"
								d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24" />
							<line v-if="!showSecret" x1="1" y1="1" x2="23" y2="23" />
							<template v-if="showSecret">
								<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
								<circle cx="12" cy="12" r="3" />
							</template>
						</svg>
					</div>
				</div>
				<div v-if="!oldUser">
					<label for="confirmSecret" class="block text-sm font-medium text-gray-700 mb-1">Confirm Secret Key</label>
					<div class="w-full flex">
						<input id="confirmSecret" @keyup="handleTextInput" :type="showConfirmSecret ? 'text' : 'password'"
							v-model="confirmSecret" placeholder="Confirm your secret key"
							class="px-4 py-2 w-96 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition" />
						<svg @click="showConfirmSecret = !showConfirmSecret" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
							fill="none" stroke="#00F376" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
							class="w-5 h-5 text-gray-400 my-auto ml-2 cursor-pointer">
							<path v-if="!showConfirmSecret"
								d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24" />
							<line v-if="!showConfirmSecret" x1="1" y1="1" x2="23" y2="23" />
							<template v-if="showConfirmSecret">
								<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
								<circle cx="12" cy="12" r="3" />
							</template>
						</svg>
					</div>
				</div>
				<div v-if="!oldUser" class="flex justify-center text-sm h-8">
					<span v-if="message" :class="message == 'Secret keys do not match.' ? 'text-red-500' : 'text-green-500'">{{
						message
					}}</span>
				</div>
				<button type="submit"
					class="w-full cursor-pointer py-3 mt-2 text-black font-semibold bg-[#00F376] hover:text-white rounded-lg shadow-lg hover:bg-black transition delay-100 ease-in-out"
					@click="handleSign()">
					{{ oldUser ? 'Sign In' : 'Sign Up' }}
				</button>
			</form>
		</div>
	</div>
</template>

<script setup lang="ts">
import type { signPayload } from '@/composables/types';
import { sign } from '~/calls/welcome';
const props = defineProps<{
	oldUser: boolean;
}>();
const message = ref('');
const emit = defineEmits(['signup', 'load', 'popup']);
const mail = ref('');
const secret = ref('');
const confirmSecret = ref('');
const appName = ref('');
const showSecret = ref(false);
const showConfirmSecret = ref(false);
const handleSign = async () => {
	let payload: signPayload = {
		mail: mail.value,
		secret: secret.value,
	}
	if (!props.oldUser) {
		if (secret.value !== confirmSecret.value) {
			return;
		} else if (mail.value === '' || secret.value === '' || appName.value === '' || confirmSecret.value === '') {
			emit('popup', 'Please fill in all fields');
			return;
		}
		payload.appName = appName.value;
	} else {
		if (mail.value === '' || secret.value === '') {
			emit('popup', 'Please fill in all fields');
			return;
		}
	}
	emit('load');
	try {
		const result = await sign(payload);
		if (result) {
			emit('signup', result);
		}
	} catch (e) {
		emit('popup', e);
	} finally {
		emit('load');
	}

}
const handleTextInput = () => {
	if (secret.value.length > 0 && confirmSecret.value.length > 0) {
		if (secret.value !== confirmSecret.value) {
			message.value = 'Secret keys do not match.';
		} else {
			message.value = 'Secret keys match. ';
		}
	}
}
</script>
