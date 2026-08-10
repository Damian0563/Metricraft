export type dashboardInitPayload = {
	appName: string;
	urls: string[];
	signedSecret: string;
	error: string;
	settings: { realtime: boolean, retention: number, enabled: Record<string, { enabled: boolean, timeframe: string }> };
}

export type config = {
	secret: string,
	httphost: string,
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

export type verifyResponse = {
	success: boolean;
	err?: string;
	status?: number;
}
