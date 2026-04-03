import type { signPayload } from '@/composables/types';
export const welcome = async (): Promise<boolean | null> => {
	const [SECRET, PORT] = useBackendUrl()
	const response = await fetch(`http://localhost:${PORT}/`, {
		headers: {
			"Authorization": SECRET,
			"Session-Token": getCookie("session-token"),
		},
		method: "GET",
	})
	const data = await response.json()
	if (data.err) {
		throw data.err
	}
	return data.exists
}

export const sign = async (payload: signPayload): Promise<boolean | null> => {
	const [SECRET, PORT] = useBackendUrl()
	const response = await fetch(`http://localhost:${PORT}/sign`, {
		headers: {
			"Authorization": SECRET,
		},
		method: "POST",
		body: JSON.stringify(payload),
	})
	const data = await response.json()
	if (data.err) {
		throw data.err
	}
	return data.token
}

function getCookie(cname: string) {
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









