

export type signPayload = {
	mail: string;
	secret: string;
	appName?: string;
}

export type dashboardInitPayload = {
	appName: string;
	signedSecret: string;
	error: string;
	settings: { realtime: boolean, retention: number, enabled: Map<string, boolean> };
}

export type config = {
	secret: string,
	httphost: string,
	wsshost: string,
	port: number,
}
