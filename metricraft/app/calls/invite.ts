import type { config, pendingUsersPayload, allowedUsersPayload } from "@/composables/types";
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

export const handlePermissionDecision = async (mail: string, action: boolean): Promise<void> => {
	const config: config = useBackendUrl()
	const headers = {
		"Authorization": config.secret,
		"Session-Token": useCookie("session-token").value || getCookie("session-token") || "",
	}
	await $fetch(`${config.httphost}/invites/handle?user=${encodeURIComponent(mail)}&action=${action}`, {
		method: "POST",
		headers,
	})
}

export const getTeamUsers = async (): Promise<{ mail: string, initials: string }[]> => {
	try {
		const config: config = useBackendUrl()
		const headers = {
			"Authorization": config.secret,
			"Session-Token": useCookie("session-token").value || getCookie("session-token") || "",
		}
		const data = await $fetch<allowedUsersPayload>(`${config.httphost}/invites/team`, {
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
