export type signPayload = {
	mail: string;
	secret: string;
	appName?: string;
}

export type dashboardInitPayload = {
	appName: string;
	urls: string[];
	signedSecret: string;
	error: string;
	settings: { realtime: boolean, retention: number, enabled: Record<string, boolean> };
}

export type config = {
	secret: string,
	httphost: string,
	wsshost: string,
	port: number,
}

export type welcomeResponse = {
	exists: boolean;
	err?: string;
}

export type signResponse = {
	token: string;
	err?: string;
}
