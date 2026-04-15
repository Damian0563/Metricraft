export const useBackendUrl = () => {
	const config = useRuntimeConfig()
	return config.public.secret
}


export const getCookie = (cname: string) => {
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

export const validateEmail = (email: string) => {
	return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
};
