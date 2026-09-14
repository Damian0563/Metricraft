export type dashboardInitPayload = {
	appName: string;
	urls: string[];
	signedSecret: string;
	error: string;
	settings: { enabled: Record<string, { enabled: boolean, timeframe: string }> };
	layout?: { name: string, span: number, height: number, custom: boolean }[];
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

export type CustomizableMetric = {
	name: string;
	timeframe: string;
	enabled: boolean;
	custom: boolean;
}

export type PreviewKind = 'map' | 'line' | 'bars' | 'donut' | 'gauge';

export type PlacedCard = {
	id: string;
	name: string;
	timeframe: string;
	kind: PreviewKind;
	custom: boolean;
	span: 1 | 2 | 3;
	height: 1 | 2 | 3;
}

export type DisplayViewCard = {
	name: string;
	span: number;
	height: number;
	custom: boolean;
}
