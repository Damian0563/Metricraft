export type dashboardInitPayload = {
	appName: string;
	urls: string[];
	signedSecret: string;
	error: string;
	settings: { retention: number, enabled: Record<string, { enabled: boolean, timeframe: string }> };
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

/* What the display-view customizer needs to list a metric in its palette:
   the enabled flag comes from dashboard settings, custom metrics are always live. */
export type CustomizableMetric = {
	name: string;
	timeframe: string;
	enabled: boolean;
	custom: boolean;
}

export type DisplayViewCard = {
	name: string;
	span: number;
	height: number;
	custom: boolean;
}
