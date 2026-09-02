export type Worker = {
	url: string;
	pollInterval: number;
	headers?: Record<string, string>;
}

export type MetricSource = 'body' | 'header' | 'query'
export type ChartType = 'line' | 'bar' | 'pie'
export type CustomMetric = {
	name: string
	method: string
	path: string
	source: MetricSource
	selector: string
	aggregation: string
	timeframe: string
	valueType: string
	applyRules: boolean
	chartType: ChartType
	lastUpdate?: string
}


export type CustomMetricResponse = {
	metrics: MetricData[];
	errors: string[];
};

export type MetricData = {
	name: string;
	metrics: GenericChartData | Map<string, number>;
	timeframe: string;
	customMetrics?: boolean;
	accumulate?: boolean;
	definition?: CustomMetric;
};


export type GenericChartDataPoint = {
	grouping: string;
	value: number;
}

export type GenericChartData = GenericChartDataPoint[];

export type Rule = {
	rule: string,
	matches: string[],
	mode: string, // "blacklist" | "grouping"
}


export type WorkerUptimeEntry = {
	stamp?: string;
	status?: number;
};

export type WorkerUptimeData = {
	entries?: WorkerUptimeEntry[];
};


export type pendingUsersPayload = {
	users: Array<{ mail: string }>;
}

export type allowedUsersPayload = {
	users: Array<{ mail: string, initials: string, status: boolean, receiveNotifications: boolean }>;
}

export type TeamUser = allowedUsersPayload["users"][number];


