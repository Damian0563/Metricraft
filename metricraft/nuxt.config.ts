// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
	compatibilityDate: "2025-07-15",
	devtools: { enabled: true },
	css: ["@/assets/css/main.css"],
	runtimeConfig: {
		public: {
			secret: process.env.SECRET || "",
			backendPort: process.env.PORT || 8000,
		},
	},
	devServer: {
		port: 3000,
	},
	vite: {
		plugins: [
			tailwindcss(),
		],
	},
});
