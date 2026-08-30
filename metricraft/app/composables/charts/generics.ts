import { emptyChart } from "@/composables/charts/charts";
import type { ChartData } from '@/composables/types/metrics'

export const genericAccumulatedChart = (
	canvas: HTMLCanvasElement,
	data: any,
): ChartData => {
	try {
		console.log(data)
		return { chart: emptyChart(canvas), additionalData: null };
	} catch (e) {
		console.error(e)
		return { chart: emptyChart(canvas), additionalData: null };
	}
}

export const genericGranularChart = (
	canvas: HTMLCanvasElement,
	data: any,
): ChartData => {
	try {
		return { chart: emptyChart(canvas), additionalData: null };
	} catch (e) {
		console.error(e)
		return { chart: emptyChart(canvas), additionalData: null };
	}
}
