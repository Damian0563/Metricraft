
export const welcome = async (): Promise<boolean | null> => {
	const [SECRET, PORT] = useBackendUrl()
	const response = await fetch(`http://localhost:${PORT}/`, {
		headers: {
			"Authorization": SECRET,
		},
		method: "GET",
	})
	const data = await response.json()
	if (data.err) {
		throw data.err
	}
	return data.exists
}













