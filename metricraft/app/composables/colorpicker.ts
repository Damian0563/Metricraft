import { oklch, formatHex } from 'culori';

const hash = (str: string): number => {
	const formatted = str.slice(str.lastIndexOf('/'))
	let h = 5381;
	for (let i = 0; i < formatted.length; i++) {
		h = ((h << 5) + h + formatted.charCodeAt(i)) >>> 0;
	}
	return h;
};

export class ColorPicker {
	private colors: Map<string, string> = new Map();
	constructor(instance: string[]) {
		this.colors = this.setColors(instance);
	}
	setColors(instance: string[]): Map<string, string> {
		let entry = new Map<string, string>();
		instance.forEach((occ: string) => {
			const slot = hash(occ)
			const hue = (slot * 137.508) % 360;
			entry.set(occ, formatHex(oklch({ mode: 'oklch', l: 0.78, c: 0.14, h: hue })));
		});
		return entry;
	}
	getColorForInstance(occ: string): string {
		return this.colors.get(occ) ?? '#000000';
	}
	destroy() {
	}
}

export function useColorPicker() {
	const urls = useState<string[]>('urls');
	const instance = useState<ColorPicker | null>('colorPicker', () => null);
	watch(
		urls,
		(newUrls) => {
			if (newUrls.length > 0) {
				instance.value = new ColorPicker(newUrls);
			}
		},
		{ immediate: true },
	);
	return instance;
}
