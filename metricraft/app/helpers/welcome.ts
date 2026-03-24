


export const welcome = async (): Promise<boolean | null> => {
	const PORT = useBackendUrl()
	const response = await fetch(`http://localhost:${PORT}/`)
	const data = await response.json()
	if (data.err) {
		throw data.err
	}
	return data.exists
}













