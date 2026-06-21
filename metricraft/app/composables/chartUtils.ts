export const createAdditionalCongestionData = (
	data: Map<string, number[]>,
	colorPicker: ColorPicker
): HTMLTableElement | null => {
	if (!data.size) return null;

	const rows = Array.from(data.entries())
		.map(([url, values]) => ({
			url,
			sum: values.reduce((s, v) => s + v, 0),
			color: colorPicker.getColorForUrl(url),
		}))
		.sort((a, b) => b.sum - a.sum);
	const table = document.createElement('table');
	table.classList.add('w-full', 'text-xs', 'border-separate', 'border-spacing-0');
	const thead = document.createElement('thead');
	thead.classList.add('top-0', 'z-10', 'bg-slate-50/95');
	const headerRow = document.createElement('tr');
	headerRow.innerHTML = `
		<th class="text-left font-semibold text-slate-400 py-1 pr-2">Endpoint</th>
		<th class="text-right font-semibold text-slate-400 py-1 pl-2 w-14">Total</th>
	`;
	thead.appendChild(headerRow);
	table.appendChild(thead);
	const tbody = document.createElement('tbody');
	rows.forEach(({ url, sum, color }, index) => {
		const row = document.createElement('tr');
		row.classList.add('transition-colors', 'hover:bg-white/70');
		if (index > 0) row.classList.add('border-t', 'border-slate-100');

		const urlCell = document.createElement('td');
		urlCell.classList.add('py-1', 'pr-2', 'min-w-0');
		const urlWrap = document.createElement('span');
		urlWrap.classList.add('inline-flex', 'items-center', 'gap-1.5', 'min-w-0', 'max-w-full');
		const dot = document.createElement('span');
		dot.classList.add('size-2', 'shrink-0', 'rounded-full', 'ring-1', 'ring-black/5');
		dot.style.backgroundColor = color;
		const label = document.createElement('span');
		label.classList.add('truncate', 'text-slate-600');
		label.textContent = url;
		label.title = url;
		urlWrap.appendChild(dot);
		urlWrap.appendChild(label);
		urlCell.appendChild(urlWrap);
		const totalCell = document.createElement('td');
		totalCell.classList.add('py-1', 'pl-2', 'text-right', 'tabular-nums', 'font-semibold', 'text-slate-700');
		totalCell.textContent = String(sum);

		row.appendChild(urlCell);
		row.appendChild(totalCell);
		tbody.appendChild(row);
	});
	table.appendChild(tbody);

	return table;
};
