import type { Chart } from "chart.js";
import type { ChoroplethChart } from 'chartjs-chart-geo';

export type Worker = {
	url: string;
	pollInterval: number;
	headers?: Record<string, string>;
}

export type HeaderRow = { key: string; value: string }

export type signPayload = {
	mail: string;
	secret: string;
	appName?: string;
}

export type ChartData = {
	chart: Chart | ChoroplethChart | null;
	additionalData: HTMLElement | null;
}

export type additionalDataHeaders = {
	h1: string;
	h2: string;
}

export type TrafficCongestionData = {
	values: CongestionEntry[];
};

export type DistributionData = {
	distribution: {
		values: Record<string, number>;
	};
};

export type ThroughputData = {
	values: ThroughputEntry[];
	computedThroughput: number;
};

export type ThroughputEntry = {
	timerange: string;
	value: number;
};

export type CongestionEntry = {
	timerange: string;
	pairing?: StringInt32Map;
};

export type StringInt32Map = {
	values?: Record<string, number>;
};

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
