export const useBackendUrl = () => {
	const config = useRuntimeConfig()
	return [config.public.secret, config.public.backendPort]
}
