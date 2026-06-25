import type { Chart } from "chart.js";

export type signPayload = {
	mail: string;
	secret: string;
	appName?: string;
}

export type ChartData = {
	chart: Chart;
	additionalData: HTMLElement | null;
}

export type dashboardInitPayload = {
	appName: string;
	urls: string[];
	signedSecret: string;
	error: string;
	settings: { realtime: boolean, retention: number, enabled: Record<string, { enabled: boolean, timeframe: string }> };
}

export type pendingUsersPayload = {
	users: Array<{ mail: string }>;
}

export type allowedUsersPayload = {
	users: Array<{ mail: string, initials: string, status: boolean }>;
}

export type TeamUser = allowedUsersPayload["users"][number];

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

export type verifyResponse = {
	success: boolean;
	err?: string;
	status?: number;
}
