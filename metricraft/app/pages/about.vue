<template>
	<div>
		<Navbar />
		<main class="px-4 sm:px-8 py-12 sm:py-16">
			<section class="max-w-5xl mx-auto text-center">
				<span
					class="inline-block px-3 py-1 rounded-full text-xs font-semibold tracking-wide uppercase bg-[#00F376]/10 text-[#00F376] border border-[#00F376]/30">
					Self-hosted log observability
				</span>
				<h1 class="mt-6 text-4xl sm:text-5xl font-bold text-white leading-tight">
					See what your API is doing,
					<span class="text-[#00F376]">without shipping your logs anywhere.</span>
				</h1>
				<p class="mt-6 text-lg text-gray-400 max-w-3xl mx-auto leading-relaxed">
					Metricraft is an analytics platform for log observability. A lightweight worker proxy captures your HTTP
					traffic, and visual dashboards turn it into latency, status code and throughput insights your whole team
					can explore.
				</p>
				<div class="mt-8 flex flex-wrap justify-center gap-3">
					<NuxtLink to="/"
						class="px-6 py-3 bg-[#00F376] text-gray-900 font-semibold rounded-lg hover:bg-[#00D96A] transition-colors duration-200">
						Get started
					</NuxtLink>
					<NuxtLink to="/manual"
						class="px-6 py-3 border border-gray-700 text-white font-semibold rounded-lg hover:border-[#00F376] hover:text-[#00F376] transition-colors duration-200">
						Read the documentation
					</NuxtLink>
				</div>
			</section>

			<section class="max-w-6xl mx-auto mt-20">
				<h2 class="text-3xl font-bold text-center text-[#00F376]">Why Metricraft</h2>
				<div class="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
					<div v-for="benefit in benefits" :key="benefit.title"
						class="bg-white rounded-xl shadow-xl p-6 border border-gray-100 hover:-translate-y-1 transition-transform duration-300">
						<div class="h-10 w-10 rounded-full bg-[#00F376]/10 flex items-center justify-center text-[#00B35C]">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
								stroke="currentColor" stroke-width="2" aria-hidden="true">
								<path stroke-linecap="round" stroke-linejoin="round" :d="benefit.icon" />
							</svg>
						</div>
						<h3 class="mt-4 text-lg font-semibold text-gray-800">{{ benefit.title }}</h3>
						<p class="mt-2 text-sm text-gray-600 leading-relaxed">{{ benefit.description }}</p>
					</div>
				</div>
			</section>

			<section class="max-w-6xl mx-auto mt-20">
				<h2 class="text-3xl font-bold text-center text-[#00F376]">What you get</h2>
				<div class="mt-8 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
					<div v-for="feature in features" :key="feature.title"
						class="rounded-xl p-6 border border-gray-800 bg-gray-900/60 hover:border-[#00F376]/60 transition-colors duration-300">
						<h3 class="text-lg font-semibold text-white">{{ feature.title }}</h3>
						<p class="mt-2 text-sm text-gray-400 leading-relaxed">{{ feature.description }}</p>
					</div>
				</div>
			</section>

			<section class="max-w-6xl mx-auto mt-20">
				<div class="bg-white rounded-xl shadow-xl p-6 sm:p-10 border border-gray-100">
					<h2 class="text-2xl font-bold text-gray-900">How it works</h2>
					<p class="mt-2 text-sm text-gray-500">From a single request to a live chart in four steps.</p>
					<ol class="mt-8 grid gap-6 md:grid-cols-4">
						<li v-for="(step, index) in flow" :key="step.title" class="relative">
							<div
								class="h-9 w-9 rounded-full bg-[#00F376] text-gray-900 font-bold flex items-center justify-center">
								{{ index + 1 }}
							</div>
							<h3 class="mt-3 font-semibold text-gray-800">{{ step.title }}</h3>
							<p class="mt-1 text-sm text-gray-600 leading-relaxed">{{ step.description }}</p>
						</li>
					</ol>
				</div>
			</section>

			<section class="max-w-6xl mx-auto mt-20">
				<h2 class="text-3xl font-bold text-center text-[#00F376]">Under the hood</h2>
				<div class="mt-8 bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden">
					<div class="overflow-x-auto">
						<table class="w-full text-left text-sm">
							<thead class="bg-gray-50 border-b border-gray-100">
								<tr>
									<th class="px-6 py-3 font-semibold text-gray-700">Component</th>
									<th class="px-6 py-3 font-semibold text-gray-700">Technology</th>
									<th class="px-6 py-3 font-semibold text-gray-700">Role</th>
								</tr>
							</thead>
							<tbody>
								<tr v-for="component in stack" :key="component.name"
									class="border-b border-gray-100 last:border-b-0 hover:bg-[#00F376]/5 transition-colors">
									<td class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap">{{ component.name }}</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span class="px-2.5 py-1 rounded-full text-xs font-medium bg-[#00F376]/20 text-green-700">
											{{ component.tech }}
										</span>
									</td>
									<td class="px-6 py-4 text-gray-600">{{ component.role }}</td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
			</section>

			<section class="max-w-4xl mx-auto mt-20 text-center">
				<h2 class="text-2xl font-bold text-white">Open source, Apache 2.0</h2>
				<p class="mt-3 text-gray-400 leading-relaxed">
					Metricraft ships as a single Docker image with the frontend, API, worker proxy, PostgreSQL and Redis
					bundled together. Run it on your own server in a few minutes.
				</p>
				<NuxtLink to="/manual#self-hosting"
					class="mt-6 inline-block px-6 py-3 bg-[#00F376] text-gray-900 font-semibold rounded-lg hover:bg-[#00D96A] transition-colors duration-200">
					Self-hosting guide
				</NuxtLink>
			</section>
		</main>
	</div>
</template>

<script setup lang="ts">
definePageMeta({
	layout: 'entry',
})

useHead({ title: 'About · Metricraft' })

const benefits = [
	{
		title: 'Self-hosted',
		description: 'Everything runs on your infrastructure. No data ever leaves your servers.',
		icon: 'M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01',
	},
	{
		title: 'Privacy-first',
		description: 'Your logs, metrics and user accounts stay in databases you own and control.',
		icon: 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z',
	},
	{
		title: 'Built for teams',
		description: 'Invite teammates and analyse traffic together, sized for small to medium engineering teams.',
		icon: 'M17 20h5v-2a4 4 0 00-3-3.87M9 20H4v-2a4 4 0 013-3.87m6-1.13a4 4 0 10-4-4 4 4 0 004 4z',
	},
	{
		title: 'Customizable',
		description: 'Define custom metrics, routing rules and dashboard layouts that match your API.',
		icon: 'M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4',
	},
]

const features = [
	{
		title: 'Real-time dashboards',
		description: 'Live charts of requests, latency and status codes, updated over WebSocket as traffic arrives.',
	},
	{
		title: 'Metricraft workers',
		description: 'Reverse proxies that record every request and poll endpoints so uptime stays tracked even at low traffic.',
	},
	{
		title: 'Overwatch',
		description: 'Extract custom metrics from request bodies, headers and query parameters, then chart them.',
	},
	{
		title: 'Grouping & blacklisting rules',
		description: 'Collapse /users/id/1, /users/id/2… into one line, or hide health checks and internal endpoints.',
	},
	{
		title: 'Custom layouts',
		description: 'Drag and drop the metrics you care about onto your own dashboard canvas.',
	},
	{
		title: 'Log retention',
		description: 'See how much storage your logs use and free it up whenever you need to.',
	},
]

const flow = [
	{ title: 'Capture', description: 'The worker proxy intercepts traffic and records headers, body, status and duration.' },
	{ title: 'Stream', description: 'Metrics stream to the Go API server over gRPC and WebSocket.' },
	{ title: 'Store', description: 'PostgreSQL keeps your logs ready for fast analytical queries.' },
	{ title: 'Visualize', description: 'The Nuxt dashboard updates in real time for everyone on your team.' },
]

const stack = [
	{ name: 'Frontend', tech: 'Nuxt 4 + Vue 3', role: 'Server-side rendered web application' },
	{ name: 'API Server', tech: 'Go', role: 'REST API and WebSocket server for real-time updates' },
	{ name: 'Worker Proxy', tech: 'Go', role: 'Reverse proxy that captures HTTP metrics and talks to the backend via gRPC' },
	{ name: 'Mail', tech: 'Go', role: 'Verification, invite, recovery and alert emails' },
	{ name: 'Metrics Store', tech: 'PostgreSQL', role: 'Log storage and analytics' },
	{ name: 'Session Cache', tech: 'Redis', role: 'Fast token validation and session management' },
	{ name: 'User Database', tech: 'PostgreSQL / Supabase', role: 'User accounts and authentication' },
]
</script>
