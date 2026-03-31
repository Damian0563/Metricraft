<template>
	<div class="flex justify-center items-center py-16 px-4">
		<div class="bg-white p-8 shadow-lg w-full max-w-lg rounded-xl">
			<div class="flex justify-center mb-6">
				<img src="/favicon.ico" alt="logo" class="w-24 h-24 rounded-lg" decoding="async" loading="lazy" />
			</div>
			<h2 class="text-2xl font-bold text-center text-gray-800 mb-6">
				{{ newUser ? 'Create Account' : 'Welcome Back' }}
			</h2>
			<form @submit.prevent class="flex flex-col gap-5">
				<div v-if="newUser">
					<label for="appName" class="block text-sm font-medium text-gray-700 mb-1">App Name</label>
					<input id="appName" type="text" placeholder="Enter your app name"
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
								<div v-if="newUser">Your email is used to send you a recovery link if you lose your secret key. We don't
									use it for anything else. For more information refer to our documentation.</div>
								<div v-else>Enter the email associated with your account to sign in.</div>
								<div class="absolute right-2 -bottom-1 w-2 h-2 bg-gray-900 rotate-45"></div>
							</div>
						</div>
					</div>
					<input id="email" type="email" placeholder="you@example.com"
						class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition mt-1" />
				</div>
				<div>
					<label for="secret" class="block text-sm font-medium text-gray-700 mb-1">Secret Key</label>
					<input id="secret" type="password" placeholder="Enter your secret key"
						class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition" />
				</div>
				<div v-if="newUser">
					<label for="confirmSecret" class="block text-sm font-medium text-gray-700 mb-1">Confirm Secret Key</label>
					<input id="confirmSecret" type="password" placeholder="Confirm your secret key"
						class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00F376] focus:border-transparent transition" />
				</div>
				<button type="submit"
					class="w-full cursor-pointer py-3 mt-2 text-black font-semibold bg-[#00F376] hover:text-white rounded-lg shadow-lg hover:bg-black transition delay-100 ease-in-out"
					@click="handleSign()">
					{{ newUser ? 'Sign Up' : 'Sign In' }}
				</button>
			</form>
		</div>
	</div>
</template>

<script setup lang="ts">
const props = defineProps<{
	newUser: boolean;
}>();

const handleSign = async () => {
	console.log("handleSign")
}

defineEmits(['signup', 'load']);
</script>
