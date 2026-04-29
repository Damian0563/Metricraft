

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


export type config = {
	secret: string,
	httphost: string,
	wsshost: string,
}
