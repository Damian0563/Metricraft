import type { Rule } from "@/composables/types";

export const getRules = async (): Promise<Rule[]> => {
	try {
		const data = await useApi()<Rule[]>(`/rules`, {
			method: "GET",
			responseType: "json",
		})
		return data;
	} catch (error) {
		console.error(error);
		return [];
	}
}

export const addRule = async (rule: Rule): Promise<{ success: boolean; err: string }> => {
	return await useApi()(`/rules`, {
		method: "POST",
		body: { rule },
	})
}

export const deleteRule = async (rule: Rule): Promise<{ success: boolean; err: string }> => {
	return await useApi()(`/rules`, {
		method: "DELETE",
		body: { rule },
	})
}
