import type { pendingUsersPayload, allowedUsersPayload, TeamUser } from "@/composables/types";
type PendingUser = { mail: string };

export const getTeamMemberStatusInfo = (status: boolean) => {
	if (status) {
		return {
			label: "Active",
			badgeClass: "w-fit shrink-0 px-3 py-1 rounded-full bg-[#00F376]/10 text-[#00A652] text-xs font-semibold uppercase tracking-wide",
			avatarClass: "h-10 w-10 shrink-0 rounded-full bg-[#00F376]/10 flex items-center justify-center text-[#00F376] font-semibold",
		};
	}
	return {
		label: "Inactive",
		badgeClass: "w-fit shrink-0 px-3 py-1 rounded-full bg-red-50 text-red-600 text-xs font-semibold uppercase tracking-wide",
		avatarClass: "h-10 w-10 shrink-0 rounded-full bg-red-50 flex items-center justify-center text-red-600 font-semibold",
	};
};

export const getPendingUsers = async (): Promise<PendingUser[]> => {
	try {
		const data = await useApi()<pendingUsersPayload>(`/invites/pending`, {
			method: "GET",
			responseType: "json",
		})
		return data.users;
	} catch (error) {
		console.error(error);
		return [];
	}
}

export const uploadUsersFromCSV = async (file: File): Promise<void> => {
	await useApi()(`/invites/batch`, {
		method: "POST",
		body: file,
	})
}

export const sendManualInvitesToUsers = async (invitees: string[]): Promise<void> => {
	await useApi()(`/invites/manual`, {
		method: "POST",
		query: { mode: "manual" },
		body: { invitees },
	})
}

export const handlePermissionDecision = async (mail: string, action: boolean): Promise<void> => {
	await useApi()(`/invites/handle`, {
		method: "POST",
		query: { user: mail, action },
	})
}

export const getTeamUsers = async (): Promise<TeamUser[]> => {
	try {
		const data = await useApi()<allowedUsersPayload>(`/invites/team`, {
			method: "GET",
			responseType: "json",
		})
		return data.users;
	} catch (error) {
		console.error(error);
		return [];
	}
}
