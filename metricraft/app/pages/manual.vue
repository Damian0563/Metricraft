<template>
	<div>
		<Navbar />
		<main class="px-4 sm:px-8 py-12">
			<header class="max-w-6xl mx-auto mb-10">
				<h1 class="text-4xl font-bold text-[#00F376]">Documentation</h1>
				<p class="mt-3 text-gray-400 max-w-3xl leading-relaxed">
					Everything you need to run Metricraft on your own server and get the most out of your dashboards.
				</p>
			</header>

			<div class="max-w-6xl mx-auto grid gap-8 lg:grid-cols-[14rem_minmax(0,1fr)]">
				<aside class="hidden lg:block">
					<nav class="sticky top-8 flex flex-col gap-6" aria-label="Documentation sections">
						<div v-for="group in toc" :key="group.title">
							<p class="text-xs font-semibold uppercase tracking-wider text-gray-500 mb-2">{{ group.title }}</p>
							<ul class="flex flex-col gap-1 border-l border-gray-800">
								<li v-for="item in group.items" :key="item.id">
									<a :href="`#${item.id}`"
										class="block -ml-px pl-4 py-1 text-sm border-l-2 transition-colors duration-200"
										:class="activeSection === item.id
											? 'border-[#00F376] text-[#00F376] font-medium'
											: 'border-transparent text-gray-400 hover:text-white'">
										{{ item.label }}
									</a>
								</li>
							</ul>
						</div>
					</nav>
				</aside>

				<div class="flex flex-col gap-6 min-w-0">
					<section id="getting-started" class="doc-section bg-white rounded-xl shadow-xl p-6 sm:p-8 border border-gray-100">
						<h2 class="text-2xl font-bold text-gray-900">Getting started</h2>
						<p class="mt-3 text-sm text-gray-600 leading-relaxed">
							Open Metricraft in your browser and choose <strong>Sign up</strong>. You will need:
						</p>
						<ul class="mt-4 grid gap-3 sm:grid-cols-3">
							<li v-for="field in signupFields" :key="field.name" class="px-4 py-3 rounded-xl border border-gray-100 bg-gray-50">
								<p class="text-sm font-medium text-gray-800">{{ field.name }}</p>
								<p class="text-xs text-gray-500 mt-1 leading-relaxed">{{ field.description }}</p>
							</li>
						</ul>
						<p class="mt-4 text-sm text-gray-600 leading-relaxed">
							After submitting, a verification code is sent to your email. Enter it to finish creating the account.
							Lost your secret key? Use <strong>Forgot password?</strong> on the sign in form to receive a recovery link.
						</p>
					</section>

					<section id="dashboard" class="doc-section bg-white rounded-xl shadow-xl p-6 sm:p-8 border border-gray-100">
						<h2 class="text-2xl font-bold text-gray-900">Dashboard</h2>
						<p class="mt-3 text-sm text-gray-600 leading-relaxed">
							The dashboard charts your traffic in real time: requests, latency and status codes per endpoint. The
							sidebar on the left gives you access to every other part of Metricraft.
						</p>
						<div class="mt-6 grid gap-4 sm:grid-cols-2">
							<div v-for="item in dashboardItems" :key="item.title" class="flex items-start gap-3">
								<div class="h-8 w-8 shrink-0 rounded-full bg-[#00F376]/10 flex items-center justify-center text-[#00B35C]">
									<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"
										aria-hidden="true">
										<path fill-rule="evenodd"
											d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
											clip-rule="evenodd" />
									</svg>
								</div>
								<div>
									<h3 class="text-sm font-semibold text-gray-800">{{ item.title }}</h3>
									<p class="text-sm text-gray-600 mt-1 leading-relaxed">{{ item.description }}</p>
								</div>
							</div>
						</div>
					</section>

					<section id="workers" class="doc-section bg-white rounded-xl shadow-xl p-6 sm:p-8 border border-gray-100">
						<h2 class="text-2xl font-bold text-gray-900">Metricraft workers</h2>
						<p class="mt-3 text-sm text-gray-600 leading-relaxed">
							Workers are lightweight reverse proxies that sit in front of your application and capture HTTP traffic
							(request method, URL, status codes, latency and headers) without changing how your app runs. Each
							worker forwards requests to your upstream service and streams metrics to the Metricraft backend.
						</p>
						<p class="mt-3 text-sm text-gray-600 leading-relaxed">
							On the <span class="doc-code">Metricraft workers</span> page, enter the upstream URL to monitor and how
							often it should be polled. Polling keeps uptime and performance metrics current even when traffic is low.
						</p>
					</section>

					<section id="overwatch" class="doc-section bg-white rounded-xl shadow-xl p-6 sm:p-8 border border-gray-100">
						<h2 class="text-2xl font-bold text-gray-900">Overwatch</h2>
						<p class="mt-3 text-sm text-gray-600 leading-relaxed">
							Overwatch lets you define <strong>custom metrics</strong> from the shape of your own API. Describe how an
							endpoint looks, then point a metric at a JSON body field, a header or a query parameter and choose how it
							should be aggregated. Metricraft charts it on your dashboards alongside everything else.
						</p>
					</section>

					<section id="rules" class="doc-section bg-white rounded-xl shadow-xl p-6 sm:p-8 border border-gray-100">
						<h2 class="text-2xl font-bold text-gray-900">Rules</h2>
						<div class="mt-4 grid gap-4 md:grid-cols-2">
							<div class="rounded-xl border border-gray-100 bg-gray-50 p-5">
								<h3 class="font-semibold text-gray-800">Grouping</h3>
								<p class="mt-2 text-sm text-gray-600 leading-relaxed">
									Collapse routes that differ only by a trailing identifier. <span class="doc-code">/users/id/1</span>
									and <span class="doc-code">/users/id/2</span> become a single
									<span class="doc-code">/users/more/</span> line. Query parameters work too:
									<span class="doc-code">/users/id?value=1</span> folds into <span class="doc-code">/users/id</span>.
								</p>
							</div>
							<div class="rounded-xl border border-gray-100 bg-gray-50 p-5">
								<h3 class="font-semibold text-gray-800">Blacklisting</h3>
								<p class="mt-2 text-sm text-gray-600 leading-relaxed">
									Hide endpoints from your dashboards, such as health checks or telemetry. Blacklisting
									<span class="doc-code">/users/id</span> also hides <span class="doc-code">/users/id/1</span> and every
									other child route.
								</p>
							</div>
						</div>
						<div class="mt-4 flex gap-3 rounded-xl border border-[#00F376]/40 bg-[#00F376]/10 p-4">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0 text-[#00B35C]" viewBox="0 0 20 20"
								fill="currentColor" aria-hidden="true">
								<path fill-rule="evenodd"
									d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
									clip-rule="evenodd" />
							</svg>
							<p class="text-sm text-gray-700 leading-relaxed">
								Always supply the full path the requests target, including the prefix, for example
								<span class="doc-code">https://api.service/users/id</span>.
							</p>
						</div>
					</section>

					<section id="team" class="doc-section bg-white rounded-xl shadow-xl p-6 sm:p-8 border border-gray-100">
						<h2 class="text-2xl font-bold text-gray-900">Team</h2>
						<p class="mt-3 text-sm text-gray-600 leading-relaxed">
							Invite teammates from the <span class="doc-code">Team</span> page by typing their emails or uploading a
							CSV file with a column of email addresses. Invited users appear under <strong>Pending verification</strong>
							until they are approved, then move to <strong>Current team</strong> with access to the project.
						</p>
					</section>

					<section id="self-hosting" class="doc-section bg-white rounded-xl shadow-xl p-6 sm:p-8 border border-gray-100">
						<h2 class="text-2xl font-bold text-gray-900">Self-hosting</h2>
						<p class="mt-3 text-sm text-gray-600 leading-relaxed">
							Metricraft ships as a single Docker image that bundles the frontend, API, worker proxy, PostgreSQL and
							Redis.
						</p>
						<div class="mt-4 rounded-xl border border-gray-100 bg-gray-50 p-5">
							<h3 class="text-sm font-semibold text-gray-800">Requirements</h3>
							<ul class="mt-2 list-disc pl-5 text-sm text-gray-600 space-y-1">
								<li>A Linux server with Docker Engine and the Compose plugin</li>
								<li>
									A Gmail account with an
									<a href="https://support.google.com/accounts/answer/185833" target="_blank" rel="noopener"
										class="text-[#00B35C] font-medium hover:underline">app password</a>,
									since sign-up sends a verification code by email
								</li>
								<li>Optionally, a domain name if you want HTTPS</li>
							</ul>
						</div>

						<ol class="mt-8 flex flex-col gap-8">
							<li v-for="(step, index) in selfHostingSteps" :key="step.title" class="flex gap-4">
								<div
									class="h-8 w-8 shrink-0 rounded-full bg-[#00F376] text-gray-900 text-sm font-bold flex items-center justify-center">
									{{ index + 1 }}
								</div>
								<div class="min-w-0 flex-1">
									<h3 class="font-semibold text-gray-800 leading-8">{{ step.title }}</h3>
									<p v-for="text in step.before" :key="text" class="mt-1 text-sm text-gray-600 leading-relaxed">
										{{ text }}
									</p>
									<div v-for="block in step.code" :key="block.label" class="mt-3 rounded-xl overflow-hidden bg-gray-900">
										<div class="flex items-center justify-between px-4 py-2 border-b border-white/10">
											<span class="text-xs font-mono text-gray-400">{{ block.label }}</span>
											<button type="button" @click="copy(block.content)"
												class="text-xs font-medium text-gray-400 hover:text-[#00F376] cursor-pointer transition-colors">
												{{ copied === block.content ? 'Copied!' : 'Copy' }}
											</button>
										</div>
										<pre class="p-4 overflow-x-auto text-sm leading-relaxed text-gray-100"><code>{{ block.content }}</code></pre>
									</div>
									<p v-for="text in step.after" :key="text" class="mt-3 text-sm text-gray-600 leading-relaxed">
										{{ text }}
									</p>
								</div>
							</li>
						</ol>
					</section>

					<section id="configuration" class="doc-section bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden">
						<div class="p-6 sm:p-8 pb-4 sm:pb-4">
							<h2 class="text-2xl font-bold text-gray-900">Configuration reference</h2>
							<p class="mt-2 text-sm text-gray-500">Environment variables read from <span class="doc-code">.env</span>.</p>
						</div>
						<div class="overflow-x-auto">
							<table class="w-full text-left text-sm">
								<thead class="bg-gray-50 border-y border-gray-100">
									<tr>
										<th class="px-6 sm:px-8 py-3 font-semibold text-gray-700">Variable</th>
										<th class="px-6 py-3 font-semibold text-gray-700">Required</th>
										<th class="px-6 sm:px-8 py-3 font-semibold text-gray-700">Description</th>
									</tr>
								</thead>
								<tbody>
									<tr v-for="variable in envVariables" :key="variable.name" class="border-b border-gray-100 last:border-b-0">
										<td class="px-6 sm:px-8 py-4 align-top whitespace-nowrap">
											<span class="doc-code">{{ variable.name }}</span>
										</td>
										<td class="px-6 py-4 align-top">
											<span class="text-xs font-medium px-3 py-1 rounded-full"
												:class="variable.required ? 'bg-[#00F376]/20 text-green-700' : 'bg-gray-200 text-gray-500'">
												{{ variable.required ? 'Yes' : 'No' }}
											</span>
										</td>
										<td class="px-6 sm:px-8 py-4 text-gray-600 leading-relaxed">{{ variable.description }}</td>
									</tr>
								</tbody>
							</table>
						</div>
					</section>

					<section id="maintenance" class="doc-section bg-white rounded-xl shadow-xl p-6 sm:p-8 border border-gray-100">
						<h2 class="text-2xl font-bold text-gray-900">Updating & backups</h2>
						<div class="mt-4 grid gap-6 md:grid-cols-2">
							<div v-for="block in maintenanceBlocks" :key="block.title" class="min-w-0">
								<h3 class="font-semibold text-gray-800">{{ block.title }}</h3>
								<p class="mt-1 text-sm text-gray-600 leading-relaxed">{{ block.description }}</p>
								<div class="mt-3 rounded-xl overflow-hidden bg-gray-900">
									<div class="flex items-center justify-between px-4 py-2 border-b border-white/10">
										<span class="text-xs font-mono text-gray-400">{{ block.label }}</span>
										<button type="button" @click="copy(block.content)"
											class="text-xs font-medium text-gray-400 hover:text-[#00F376] cursor-pointer transition-colors">
											{{ copied === block.content ? 'Copied!' : 'Copy' }}
										</button>
									</div>
									<pre class="p-4 overflow-x-auto text-sm leading-relaxed text-gray-100"><code>{{ block.content }}</code></pre>
								</div>
							</div>
						</div>
					</section>
				</div>
			</div>
		</main>
	</div>
</template>

<script setup lang="ts">
definePageMeta({
	layout: 'entry',
})

useHead({ title: 'Documentation · Metricraft' })

const toc = [
	{
		title: 'Using Metricraft',
		items: [
			{ id: 'getting-started', label: 'Getting started' },
			{ id: 'dashboard', label: 'Dashboard' },
			{ id: 'workers', label: 'Workers' },
			{ id: 'overwatch', label: 'Overwatch' },
			{ id: 'rules', label: 'Rules' },
			{ id: 'team', label: 'Team' },
		],
	},
	{
		title: 'Self-hosting',
		items: [
			{ id: 'self-hosting', label: 'Installation' },
			{ id: 'configuration', label: 'Configuration' },
			{ id: 'maintenance', label: 'Updating & backups' },
		],
	},
]

const signupFields = [
	{ name: 'App name', description: 'The application you are monitoring. Metrics and accounts are grouped under it.' },
	{ name: 'Email', description: 'Used only for verification and recovery links if you lose your secret key.' },
	{ name: 'Secret key', description: 'Your password. Pick a strong one and confirm it.' },
]

const dashboardItems = [
	{ title: 'Custom layouts', description: 'Open Settings → Customize Dashboard View and drag metrics onto your own canvas.' },
	{ title: 'Derived metrics', description: 'Toggle additional metrics computed from your logs, then apply the changes.' },
	{ title: 'Log retention', description: 'Check current log storage in Settings and delete logs to free up space.' },
	{ title: 'Rules', description: 'Jump to grouping and blacklisting rules straight from Settings.' },
]

const selfHostingSteps = [
	{
		title: 'Create the configuration',
		before: ['Create a directory for Metricraft, then add a .env file with your values and a compose.yaml next to it.'],
		code: [
			{ label: 'terminal', content: 'mkdir metricraft && cd metricraft' },
			{
				label: '.env',
				content: `# Name of the app you are monitoring
APPNAME=my-app
# Port your app listens on; the worker proxy forwards traffic here
DEST_PORT=3000
# API URL as the browser reaches it (port 8080, or your API domain)
NUXT_PUBLIC_HTTPHOST=http://localhost:8080
# Generate with: openssl rand -hex 32
SECRET=
# Users database. The bundled PostgreSQL works; any PostgreSQL (e.g. Supabase) does too
DATABASE_USERS=postgresql://postgres:password@127.0.0.1:5432/postgres?sslmode=disable
# Gmail account that sends verification, invite and alert emails
GOOGLE_MAIL_ADDRESS=you@gmail.com
GOOGLE_APP_PASSWORD=`,
			},
			{
				label: 'compose.yaml',
				content: `services:
  metricraft:
    image: damianek952/metricraft:latest   # pin a release tag in production
    restart: unless-stopped
    stop_grace_period: 30s
    env_file: .env
    ports:
      - "127.0.0.1:8000:8000"   # dashboard
      - "127.0.0.1:8080:8080"   # API
      - "127.0.0.1:8081:8081"   # worker proxy
    volumes:
      - metricraft-db:/var/lib/postgresql/data

volumes:
  metricraft-db:`,
			},
		],
		after: [
			'The ports only listen on 127.0.0.1, so they are reachable from the server itself or through a reverse proxy. To expose them directly, remove the 127.0.0.1: prefix.',
		],
	},
	{
		title: 'Start Metricraft',
		before: [],
		code: [{ label: 'terminal', content: 'docker compose up -d\ndocker compose logs -f metricraft' }],
		after: [],
	},
	{
		title: 'Create the users tables',
		before: [
			'Metricraft creates its log tables automatically, but the account tables are created once in the database DATABASE_USERS points to. For the bundled PostgreSQL:',
		],
		code: [
			{
				label: 'terminal',
				content: `docker compose exec -T metricraft psql -U postgres -d postgres <<'SQL'
CREATE TABLE IF NOT EXISTS public.users (
  created_at    timestamptz NOT NULL DEFAULT now(),
  app_name      text,
  mail          text PRIMARY KEY,
  secret        text NOT NULL,
  uuid          uuid,
  allowed_users jsonb DEFAULT '[]'::jsonb,
  pending_users jsonb DEFAULT '[]'::jsonb,
  owner         boolean DEFAULT false
);
CREATE TABLE IF NOT EXISTS public.workers (
  id         bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT now(),
  app_name   text,
  workers    json DEFAULT '[]'::json
);
SQL
docker compose restart metricraft`,
			},
		],
		after: [
			'If you use Supabase or another external database, run the same SQL there. Then open http://localhost:8000 and sign up.',
		],
	},
	{
		title: 'Expose it on a domain',
		before: ['Put a reverse proxy in front of Metricraft for HTTPS. For example, with Caddy, which issues certificates automatically:'],
		code: [
			{
				label: 'Caddyfile',
				content: `metrics.example.com {
	reverse_proxy 127.0.0.1:8000
}

api.metrics.example.com {
	reverse_proxy 127.0.0.1:8080
}

app.example.com {
	reverse_proxy 127.0.0.1:8081
}`,
			},
		],
		after: [
			'Then set NUXT_PUBLIC_HTTPHOST=https://api.metrics.example.com in .env and run docker compose up -d. No rebuild is needed.',
		],
	},
	{
		title: "Route your app's traffic through the worker proxy",
		before: [
			'The worker proxy on :8081 records each request and forwards it to http://<request host>:<DEST_PORT>. In the example above, a request to app.example.com is forwarded to http://app.example.com:3000.',
			'Your app must be reachable on DEST_PORT at that hostname from inside the container, and the reverse proxy must pass the original Host header without a port (Caddy does by default; for nginx use proxy_set_header Host $host;).',
		],
		code: [],
		after: [],
	},
]

const envVariables = [
	{ name: 'APPNAME', required: true, description: 'Name of the monitored app. Metrics and accounts are grouped under it.' },
	{ name: 'SECRET', required: true, description: 'Token shared by the frontend and API. It is visible to the browser, so treat it as an API key, not a password.' },
	{ name: 'DATABASE_USERS', required: true, description: 'PostgreSQL connection string for the users database.' },
	{ name: 'NUXT_PUBLIC_HTTPHOST', required: true, description: 'Public API URL the browser uses. Must not be an internal Docker hostname.' },
	{ name: 'GOOGLE_MAIL_ADDRESS', required: true, description: 'Gmail address that sends verification, invite, recovery and alert emails.' },
	{ name: 'GOOGLE_APP_PASSWORD', required: true, description: 'App password for GOOGLE_MAIL_ADDRESS.' },
	{ name: 'DEST_PORT', required: false, description: 'Port the worker proxy forwards to. Default: 3000.' },
]

const maintenanceBlocks = [
	{
		title: 'Updating',
		description: 'Logs and metrics live in the metricraft-db volume and survive updates. Read the release notes before upgrading across a PostgreSQL major version.',
		label: 'terminal',
		content: 'docker compose pull\ndocker compose up -d',
	},
	{
		title: 'Backup and restore',
		description: 'Dump the database to a file, and restore it into a fresh metricraft-db volume.',
		label: 'terminal',
		content: `docker compose exec -T metricraft pg_dump -U postgres postgres > metricraft-backup.sql
docker compose exec -T metricraft psql -U postgres -d postgres < metricraft-backup.sql`,
	},
]

const copied = ref('')
let copiedTimeout: ReturnType<typeof setTimeout> | undefined

const copy = async (content: string) => {
	try {
		await navigator.clipboard.writeText(content)
		copied.value = content
		clearTimeout(copiedTimeout)
		copiedTimeout = setTimeout(() => (copied.value = ''), 2000)
	} catch (e) {
		console.log(e)
	}
}

const activeSection = ref('getting-started')
let observer: IntersectionObserver | undefined

onMounted(() => {
	observer = new IntersectionObserver((entries) => {
		const visible = entries.filter((entry) => entry.isIntersecting)
		if (visible.length) {
			activeSection.value = visible[0]!.target.id
		}
	}, { rootMargin: '-20% 0px -70% 0px' })
	document.querySelectorAll('.doc-section').forEach((section) => observer!.observe(section))
})

onBeforeUnmount(() => {
	observer?.disconnect()
	clearTimeout(copiedTimeout)
})
</script>

<style scoped>
.doc-section {
	scroll-margin-top: 2rem;
}

.doc-code {
	font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
	font-size: 0.8125rem;
	color: #1f2937;
	background-color: #f3f4f6;
	padding: 0.1rem 0.35rem;
	border-radius: 0.25rem;
}
</style>
