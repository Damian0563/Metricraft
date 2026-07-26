import { Chart } from "chart.js";
import { emptyChart } from "~/composables/charts/charts";
import type { WorkerUptimeData, WorkerUptimeEntry } from '@/composables/types/additional'

export const createWorkerUptimeChart = (
	canvas: HTMLCanvasElement,
	data: WorkerUptimeData,
): Chart => {
	const upColor = '#00F376';
	const downColor = '#DC2626';
	const noColor = '#475569';
	try {
		const entries: WorkerUptimeEntry[] = data?.entries ?? [];
		if (!entries.length) throw new Error('Data is empty');
		const dates: Date[] = entries.map((entry) => new Date(entry.stamp ?? 0));
		const labels: string[] = dates.map((date) =>
			date.toLocaleString(undefined, { timeZone: "UTC", month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
		);
		const statuses: number[] = entries.map((entry) => entry.status!);
		const colors: string[] = statuses.map((up) => (up === 1 ? upColor : up === 0 ? downColor : noColor));
		const upCount: number = statuses.filter(num => Math.abs(num) === 1).length;
		const uptimePct: number = statuses.length ? (upCount / statuses.length) * 100 : 0;
		return new Chart(canvas, {
			type: 'bar',
			data: {
				labels,
				datasets: [{
					label: 'Status',
					data: statuses.map(() => 1),
					backgroundColor: colors,
					hoverBackgroundColor: colors,
					borderWidth: 0,
					borderRadius: 2,
					borderSkipped: false,
					categoryPercentage: 1,
					barPercentage: 0.92,
				}],
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				devicePixelRatio: Math.ceil(window.devicePixelRatio || 1),
				animation: { duration: 250 },
				interaction: { mode: 'index', intersect: false },
				layout: { padding: { right: 8, left: 2, top: 4, bottom: 2 } },
				scales: {
					x: {
						grid: { display: false },
						border: { display: false },
						ticks: {
							autoSkip: false,
							maxRotation: 0,
							padding: 6,
							color: '#475569',
							font: { weight: 'bold', size: 11 },
							callback: (_value: string | number, index: number): string =>
								index === 0 || index === labels.length - 1 ? (labels[index] ?? '') : '',
						},
					},
					y: {
						display: false,
						beginAtZero: true,
						max: 1,
					},
				},
				plugins: {
					legend: { display: false },
					title: {
						display: true,
						text: `Uptime: ${uptimePct.toFixed(1)}%, last 31 days`,
						color: '$fffff',
						font: { weight: 'bold', size: 18 },
						padding: { bottom: 8 },
					},
					tooltip: {
						enabled: true,
						displayColors: true,
						padding: 10,
						titleFont: { weight: 'bold', size: 13 },
						bodyFont: { size: 12 },
						callbacks: {
							title: (items) => {
								const index = items[0]?.dataIndex ?? 0;
								return dates[index]?.toLocaleString(undefined, { timeZone: "UTC", month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }) ?? '';
							},
							label: (item) => (statuses[item.dataIndex] === 1 ? 'Status: Up' : statuses[item.dataIndex] === 0 ? 'Status: Down' : 'Status: No data'),
							labelColor: (item) => {
								const color = statuses[item.dataIndex] === 1 ? upColor : statuses[item.dataIndex] === 0 ? downColor : noColor;
								return { borderColor: color, backgroundColor: color };
							},
						},
					},
				},
			},
		});
	} catch (_) {
		return emptyChart(canvas);
	}
};
