import type { signPayload, config, welcomeResponse, signResponse } from '@/composables/types';
import { getCookie } from "@/composables/helpers";
export const welcome = async (): Promise<boolean> => {
	try {
		const config: config = useBackendUrl()
		const headers = {
			"Authorization": config.secret,
			"Session-Token": useCookie("session-token").value || getCookie("session-token") || "",
		}
		const response = await $fetch<welcomeResponse>(`${config.httphost}/welcome`, {
			method: "GET", headers
		})
		if (response.err) throw (response.err)
		return Boolean(response.exists)
	} catch (e) {
		return false
	}
}

export const sign = async (payload: signPayload): Promise<string | null> => {
	const config: config = useBackendUrl()
	const headers = {
		"Authorization": config.secret,
		"Content-Type": "application/json",
	}
	try {
		let data: string = await $fetch(`${config.httphost}/sign`, {
			method: 'POST',
			headers,
			body: payload,
		})
		const object: signResponse = JSON.parse(data)
		if (object.err) {
			throw object.err
		}
		return object.token
	} catch (e: any) {
		if (e.status === 429) {
			throw "Too many requests, try again later."
		}
		throw "Something went wrong, Check your internet connection and try again."
	}
}

export const verify = async (mail: string, code: string): Promise<boolean> => {
	const config: config = useBackendUrl()
	const headers = {
		"Authorization": config.secret,
		"Content-Type": "application/json",
	}
	try {
		let data: string = await $fetch(`${config.httphost}/verify`, {
			method: 'POST',
			headers,
			body: JSON.stringify({ mail, code }),
		})
		const object: verifyResponse = JSON.parse(data)
		console.log("Verified: ", object)
		if (object.err) {
			throw object.err
		}
		return object.success
	} catch (e: any) {
		if (e.status === 429) {
			throw "Too many requests, try again later."
		}
		throw "Something went wrong, Check your internet connection and try again."
	}
}

export const sendVerification = async (mail: string): Promise<boolean> => {
	const config: config = useBackendUrl()
	const headers = {
		"Authorization": config.secret,
		"Content-Type": "application/json",
	}
	try {
		let data: string = await $fetch(`${config.httphost}/verify/send`, {
			method: 'POST',
			headers,
			body: mail,
		})
		const object: verifyResponse = JSON.parse(data)
		console.log("Sent: ", object)
		if (object.err) {
			throw object.err
		}
		return object.success
	} catch (e: any) {
		if (e.status === 429) {
			throw "Too many requests, try again later."
		}
		throw "Something went wrong, Check your internet connection and try again."
	}
}









