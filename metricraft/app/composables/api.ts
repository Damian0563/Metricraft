import type { FetchOptions } from "ofetch";
import { getCookie } from "@/composables/helpers";

type ApiFetch = <T = unknown>(request: string, opts?: FetchOptions) => Promise<T>;
export const useApi = (): ApiFetch => {
	const config = useRuntimeConfig();
	const sessionToken = useCookie<string>("session-token");
	return $fetch.create({
		baseURL: config.public.httphost as string,
		onRequest({ options }) {
			const headers = new Headers(options.headers);
			headers.set("Authorization", config.public.secret as string);
			const token = sessionToken.value || getCookie("session-token") || "";
			if (token) headers.set("Session-Token", token);
			options.headers = headers;
		},
	}) as ApiFetch;
};
