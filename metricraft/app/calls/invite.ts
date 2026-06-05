import type { config, pendingUsersPayload } from "@/composables/types";
import { getCookie } from "@/composables/helpers";
type PendingUser = { mail: string };

export const getPendingUsers = async (): Promise<PendingUser[]> => {
	try {
		const config: config = useBackendUrl()
		const headers = {
			"Authorization": config.secret,
			"Session-Token": useCookie("session-token").value || getCookie("session-token") || "",
		}
		const data = await $fetch<pendingUsersPayload>(`${config.httphost}/invites/pending`, {
			method: "GET",
			headers,
			responseType: "json",
		})
		return data.users;
	} catch (error) {
		console.error(error);
		return [];
	}
}
