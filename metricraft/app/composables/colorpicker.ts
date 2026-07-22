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
	private urls_color: Map<string, string> = new Map();
	constructor(urls: string[]) {
		this.urls_color = this.setColors(urls);
	}
	setColors(urls: string[]): Map<string, string> {
		let entry = new Map<string, string>();
		urls.forEach((url: string) => {
			const slot = hash(url)
			const hue = (slot * 137.508) % 360;
			entry.set(url, formatHex(oklch({ mode: 'oklch', l: 0.78, c: 0.14, h: hue })));
		});
		return entry;
	}
	getColorForUrl(url: string): string {
		return this.urls_color.get(url) ?? '#000000';
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
