
export const toggleRealtime = async (enabled: boolean) => {
	try {
		const SECRET = useBackendUrl()
		const response = await fetch(`http://localhost:8080/settings/realtime`, {
			method: "POST",
			headers: {
				"Authorization": SECRET,
				"Session-Token": getCookie("session-token"),
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ enabled }),
		});
		console.log(response);
		console.log(enabled);
		if (!response.ok) throw new Error("Failed to toggle realtime");
	} catch (error) {
		console.error(error);
	}
}
