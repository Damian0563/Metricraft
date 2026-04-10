
export const handleMessage = (event: MessageEvent) => {
	const data = JSON.parse(event.data)
	console.log(data)
}
