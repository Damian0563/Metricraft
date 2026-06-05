import type { signPayload, config, welcomeResponse, signResponse, verifyResponse } from '@/composables/types';
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

export const verify = async (mail: string, code: string, appName: string): Promise<verifyResponse> => {
	const config: config = useBackendUrl()
	const headers = {
		"Authorization": config.secret,
		"Content-Type": "application/json",
	}
	try {
		await $fetch(`${config.httphost}/verify/check`, {
			method: 'POST',
			headers,
			body: JSON.stringify({ mail: mail, code: code, appName: appName }),
		})
		const res: verifyResponse = {
			success: true,
		}
		return res
	} catch (e: any) {
		let error: string
		if (e.status === 429) {
			error = "Too many requests, try again later."
		} else if (e.status === 400) {
			error = "Verification code is invalid or has expired."
		} else if (e.status === 401) {
			error = "You don't have permission to access this app, please contact the owner. Once you have access, you will be able to sign up."
		} else if (e.status === 403) {
			error = "Invalid app name. Please contact the administrator."
		} else {
			error = "Something went wrong, Check your internet connection and try again."
		}
		const res: verifyResponse = {
			success: false,
			err: error,
			status: e.status,
		}
		return res
	}
}

export const sendVerification = async (mail: string): Promise<verifyResponse> => {
	const config: config = useBackendUrl()
	const headers = {
		"Authorization": config.secret,
		"Content-Type": "application/json",
	}
	try {
		await $fetch(`${config.httphost}/verify/send`, {
			method: 'POST',
			headers,
			body: JSON.stringify({ "mail": mail }),
		})
		const res: verifyResponse = {
			success: true,
		}
		return res
	} catch (e: any) {
		let error: string
		if (e.status === 429) {
			error = "Too many requests, try again later."
		} else if (e.status === 400) {
			error = "Email already exists."
		} else {
			error = "Something went wrong, Check your internet connection and try again."
		}
		const res: verifyResponse = {
			success: false,
			err: error,
		}
		return res
	}
}









