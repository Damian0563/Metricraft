import type { additionalDataHeaders } from "~/composables/types";
import type { ColorPicker } from "~/composables/colorpicker";

export const createAdditionalData = (
	data: Map<string, number>,
	headers: additionalDataHeaders,
	colorPicker: ColorPicker | null = null
): HTMLTableElement | null => {
	if (!data.size) return null;
	const rows = Array.from(data.entries())
		.map(([url, values]) => ({
			url,
			sum: values,
			color: colorPicker?.getColorForUrl(url),
		}))
		.sort((a, b) => b.sum - a.sum);
	const table = document.createElement('table');
	table.classList.add('w-full', 'text-xs', 'table-fixed', 'border-separate', 'border-spacing-0');
	const colgroup = document.createElement('colgroup');
	const labelCol = document.createElement('col');
	labelCol.classList.add('w-[58%]');
	const valueCol = document.createElement('col');
	valueCol.classList.add('w-[42%]');
	colgroup.append(labelCol, valueCol);
	table.appendChild(colgroup);
	const thead = document.createElement('thead');
	thead.classList.add('top-0', 'z-10', 'bg-slate-50/50');
	const headerRow = document.createElement('tr');
	headerRow.innerHTML = `
		<th class="text-left font-semibold text-slate-400 py-1.5 px-3 whitespace-nowrap">${headers.h1}</th>
		<th class="text-right font-semibold text-slate-400 py-1.5 px-3 whitespace-nowrap">${headers.h2}</th>
	`;
	thead.appendChild(headerRow);
	table.appendChild(thead);
	const tbody = document.createElement('tbody');
	rows.forEach(({ url, sum, color }, index) => {
		const row = document.createElement('tr');
		row.classList.add('transition-colors', 'hover:bg-white/70');
		if (index > 0) row.classList.add('border-t', 'border-slate-100');
		const urlCell = document.createElement('td');
		urlCell.classList.add('py-1.5', 'px-3', 'min-w-0');
		const urlWrap = document.createElement('span');
		urlWrap.classList.add('inline-flex', 'items-center', 'gap-1.5', 'min-w-0', 'max-w-full');
		if (colorPicker && color) {
			const dot = document.createElement('span');
			dot.classList.add('size-2', 'shrink-0', 'rounded-full', 'ring-1', 'ring-black/5');
			dot.style.backgroundColor = color;
			urlWrap.appendChild(dot);
		}
		const label = document.createElement('span');
		label.classList.add('truncate', 'text-slate-600');
		label.textContent = url;
		label.title = url;
		urlWrap.appendChild(label);
		urlCell.appendChild(urlWrap);
		const totalCell = document.createElement('td');
		totalCell.classList.add('py-1.5', 'px-3', 'text-right', 'tabular-nums', 'font-semibold', 'text-slate-700', 'whitespace-nowrap');
		totalCell.textContent = String(sum);
		row.appendChild(urlCell);
		row.appendChild(totalCell);
		tbody.appendChild(row);
	});
	table.appendChild(tbody);
	return table;
};
