import type { Rule } from '@/composables/types/additional'

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

export const addRule = async (rule: Rule): Promise<void> => {
	return await useApi()(`/rules/add`, {
		method: "POST",
		body: rule,
	})
}

export const deleteRule = async (rule: Rule): Promise<void> => {
	return await useApi()(`/rules/delete`, {
		method: "DELETE",
		body: rule,
	})
}
