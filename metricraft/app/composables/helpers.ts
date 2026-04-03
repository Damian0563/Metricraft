export const useBackendUrl = () => {
	const config = useRuntimeConfig()
	return [config.public.secret, config.public.backendPort]
}


export const getCookie = (cname: string) => {
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
}
