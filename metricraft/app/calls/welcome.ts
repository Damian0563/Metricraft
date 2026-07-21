import type { signPayload, welcomeResponse, signResponse, verifyResponse } from '@/composables/types';

export const welcome = async (): Promise<boolean> => {
	try {
		const response = await useApi()<welcomeResponse>(`/welcome`, {
			method: "GET",
		})
		const result: welcomeResponse = typeof response === 'string' ? JSON.parse(response) : response
		if (result.err) return false
		return Boolean(result.exists)
	} catch (e) {
		return false
	}
}

export const sign = async (payload: signPayload): Promise<string | null> => {
	try {
		const data: string | signResponse = await useApi()(`/sign`, {
			method: 'POST',
			body: payload,
		})
		const object: signResponse = typeof data === 'string' ? JSON.parse(data) : data
		if (object.err) {
			throw object.err
		}
		return object.token
	} catch (e: any) {
		if (e.status === 429) {
			throw "Too many requests, try again later."
		}
		if (typeof e === 'string') {
			throw e
		}
		if (e.data?.err) {
			throw e.data.err
		}
		if (e.status === 401) {
			throw "Invalid credentials."
		}
		throw "Something went wrong, Check your internet connection and try again."
	}
}

export const verify = async (mail: string, code: string, appName: string): Promise<verifyResponse> => {
	try {
		await useApi()(`/verify/check`, {
			method: 'POST',
			body: { mail, code, appName },
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

export const sendRecovery = async (mail: string): Promise<void> => {
	try {
		await useApi()(`/recovery/send`, {
			method: 'POST',
			body: { mail },
		})
	} catch (e: any) {
		if (e.status === 429) {
			throw "Too many requests, try again later."
		}
		throw "Something went wrong, Check your internet connection and try again."
	}
}

export const checkRecovery = async (id: string, password: string): Promise<void> => {
	try {
		await useApi()(`/recovery/check`, {
			method: 'POST',
			query: { id },
			body: { password },
		})
	} catch (e: any) {
		if (e.status === 429) {
			throw "Too many requests, try again later."
		}
		else if (e.status === 400) {
			throw "Bad request, try again."
		} else if (e.status === 403 || e.status === 401) {
			throw "Recovery link is invalid or has expired."
		}
		throw "Something went wrong, Check your internet connection and try again."
	}
}

export const sendVerification = async (mail: string): Promise<verifyResponse> => {
	try {
		await useApi()(`/verify/send`, {
			method: 'POST',
			body: { mail },
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
