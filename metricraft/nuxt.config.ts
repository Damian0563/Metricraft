// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";
import { fileURLToPath } from "node:url";
import fs from "node:fs";
import path from "node:path";

const getPort = (): number | undefined => {
	try {
		const configPath = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..", "config.json");
		const config = JSON.parse(fs.readFileSync(configPath, "utf8"));
		return config.port;
	} catch (e) {
		console.log(e);
	}
};


export default defineNuxtConfig({
	compatibilityDate: "2025-07-15",
	devtools: { enabled: true },
	css: ["@/assets/css/main.css"],
	runtimeConfig: {
		public: {
			secret: process.env.SECRET || "",
			backendPort: 8080,
		},
	},
	devServer: {
		port: getPort(),
	},
	vite: {
		plugins: [
			tailwindcss(),
		],
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
	},
});
