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

export const getCookie = (cname: string) => {
	if (process.server) return ""
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

export const updateCookie = (cvalue: string) => {
	document.cookie = `session-token=${cvalue};path=/;expires=${new Date(Date.now() + 3600000 * 24).toUTCString()};SameSite=None;Secure`;
}

export const validateEmail = (email: string): boolean => {
	return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
};

export const evaluatePasswordStrength = (password: string): number => {
	const passwordScore = zxcvbn(password);
	return passwordScore.score;
};

export const invalidateCookie = () => {
	document.cookie = "session-token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
}
