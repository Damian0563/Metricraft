import tailwindcss from "@tailwindcss/vite";
import type { config } from "./app/composables/types";

const getConfig = (): config => {
	try {
		let mode = import.meta.env.MODE;
		if (mode === undefined) {
			throw new Error("no config file");
		}
		let config: config = { "secret": import.meta.env.SECRET || "", "httphost": "", "wsshost": "", "port": 0 };
		const httpHost = import.meta.env.NUXT_PUBLIC_HTTPHOST;
		const wsHost = import.meta.env.NUXT_PUBLIC_WSSHOST;
		if (httpHost) {
			config.httphost = httpHost;
			config.wsshost = wsHost || httpHost.replace(/^http/, "ws");
		} else if (mode === "local") {
			config.httphost = "http://localhost:8080";
			config.wsshost = "ws://localhost:8080";
		}
		else {
			throw new Error("invalid mode");
		}
		return config;
	} catch (e) {
		console.log(e);
		return { "port": 0, "httphost": "", "wsshost": "", "secret": "" };
	}
};


export default defineNuxtConfig({
	compatibilityDate: "2025-07-15",
	devtools: { enabled: true },
	css: ["@/assets/css/main.css"],
	runtimeConfig: {
		public: {
			backendPort: 8080,
			secret: getConfig().secret,
			wsshost: getConfig().wsshost,
			httphost: getConfig().httphost,
		},
		secret: getConfig().secret,
	},
	devServer: {
		port: 8000,
	},
	vite: {
		plugins: [
			tailwindcss(),
		],
		optimizeDeps: {
			include: [
				'chart.js',
				'culori',
			]
		}
	},
	app: {
		head: {
			title: "Metricraft",
			meta: [
				{ charset: "utf-8" },
				{ name: "viewport", content: "width=device-width, initial-scale=1" },
			],
			link: [
				{ rel: "icon", type: "image/ico", href: "/favicon.ico" },
			],
		},
		pageTransition: { name: 'page', mode: 'out-in' },
	},
});
