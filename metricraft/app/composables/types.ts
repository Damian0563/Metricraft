

export type signPayload = {
	mail: string;
	secret: string;
	appName?: string;
}

export type dashboardInitPayload = {
	appName: string;
	signedSecret: string;
	error: string;
	settings: { realtime: boolean };
}
