// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";
import { config } from "./composables/types";

const getConfig = (): config => {
	try {
		let mode = import.meta.env.MODE;
		if (mode === undefined) {
			throw new Error("no config file");
		}
		let config: config = { "secret": import.meta.env.SECRET || "", "httphost": "", "wsshost": "", "port": 0 };
		const domain = import.meta.env.DOMAIN;
		if (domain) {
			// DOMAIN is the public backend base URL the browser uses, e.g.
			// "https://metrics.example.com" or "http://203.0.113.10:8080".
			config.httphost = domain;
			config.wsshost = domain.replace(/^http/, "ws");
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
			secret: getConfig().secret,
			backendPort: 8080,
			wsshost: getConfig().wsshost,
			httphost: getConfig().httphost,
		},
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
