import { zxcvbn } from "@zxcvbn-ts/core"

export const useBackendUrl = (): config => {
	const config = useRuntimeConfig()
	return {
		secret: config.public.secret,
		httphost: config.public.httphost,
		wsshost: config.public.wsshost,
		port: config.public.backendPort,
	}
}

export const truncateUrl = (url: string, max = 20): string => {
	if (url.length <= max) return url;
	return `${url.slice(0, max - 1)}…`;
};
export const getCookie = (cname: string): string => {
	if (import.meta.server) return ""
	try {
		let name = cname + "=";
		let decodedCookie = decodeURIComponent(document.cookie);
		let ca = decodedCookie.split(';');
		for (let i = 0; i < ca.length; i++) {
			let c = ca[i];
			if (!c || !c.trim()) {
				continue;
			}
			while (c.charAt(0) == ' ') {
				c = c.substring(1);
			}
			if (c.indexOf(name) == 0) {
				return c.substring(name.length, c.length);
			}
		}
		return "";
	} catch (e) {
		console.log(e)
		return ""
	}
}

export const updateCookie = (cvalue: string): void => {
	document.cookie = `session-token=${cvalue};path=/;expires=${new Date(Date.now() + 3600000 * 24).toUTCString()};SameSite=None;Secure`;
}

export const validateEmail = (email: string): boolean => {
	return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
};

export const evaluatePasswordStrength = (password: string): number => {
	const passwordScore = zxcvbn(password);
	return passwordScore.score;
};

export const invalidateCookie = (): void => {
	document.cookie = "session-token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
}

export const parseApiError = (error: unknown, fallback: string): string => {
	if (!error || typeof error !== 'object') return fallback

	const { data } = error as { data?: unknown }
	if (data && typeof data === 'object') {
		const record = data as Record<string, unknown>
		if (typeof record.err === 'string' && record.err.trim()) return record.err.trim()
	}

	if (typeof data === 'string') {
		const trimmed = data.trim()
		const jsonStart = trimmed.indexOf('{')
		if (jsonStart !== -1) {
			try {
				const parsed = JSON.parse(trimmed.slice(jsonStart)) as Record<string, unknown>
				if (typeof parsed.err === 'string' && parsed.err.trim()) return parsed.err.trim()
			} catch {
				// fall through to plain-text handling
			}
		}
		const firstLine = trimmed.split('\n').map((line) => line.trim()).find(Boolean)
		if (firstLine && !firstLine.startsWith('{')) return firstLine
	}

	return fallback
}
