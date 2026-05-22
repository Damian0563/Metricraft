import { Chart } from "chart.js";

type StringInt32Map = {
	values?: Record<string, number>;
};
type CongestionEntry = {
	timerange: string;
	pairing?: StringInt32Map;
};
export type TrafficCongestionData = {
	values: CongestionEntry[];
};

export const createTrafficCongestionTrends = (
	canvas: HTMLCanvasElement,
	data: TrafficCongestionData,
	colorPicker: ColorPicker
): Chart => {
	const getUrlCounts = (pairing: StringInt32Map | undefined): Record<string, number> =>
		pairing?.values ?? {};
	try {
		if (!data?.values?.length) throw new Error('Data is empty');
		const labels: string[] = [];
		const urlDataMap = new Map<string, number[]>();
		const totalPoints = data.values.length;
		data.values.forEach((entry: CongestionEntry, pointIndex: number) => {
			labels.push(entry.timerange);
			for (const [url, count] of Object.entries(getUrlCounts(entry.pairing))) {
				if (!urlDataMap.has(url)) {
					urlDataMap.set(url, new Array<number>(totalPoints).fill(0));
				}
				urlDataMap.get(url)![pointIndex] = count;
			}
		});
		console.log(colorPicker)
		const datasets = Array.from(urlDataMap.entries()).map(([url, dataArray], _) => {
			return {
				label: url,
				data: dataArray,
				borderWidth: 1,
				backgroundColor: colorPicker.getColorForUrl(url),
			};
		});
		return new Chart(canvas, {
			type: 'bar',
			data: {
				labels,
				datasets
			},
			options: {
				interaction: {
					mode: 'nearest',
					intersect: true
				},
				responsive: true,
				maintainAspectRatio: false,
				scales: {
					x: { stacked: true },
					y: { beginAtZero: true, stacked: true }
				},
				plugins: {
					legend: {
						display: true,
						position: 'top'
					},
					tooltip: {
						enabled: true,
					}
				},
			}
		});
	} catch (e) {
		console.error(e);
		return new Chart(canvas, {
			type: 'bar',
			data: { labels: [], datasets: [] }
		});
	}
}
