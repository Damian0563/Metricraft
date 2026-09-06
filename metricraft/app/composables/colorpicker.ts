import { oklch, formatHex } from 'culori';

const GOLDEN_ANGLE = 137.508;
const LIGHTNESS = [0.78, 0.66, 0.86];
const CHROMA = [0.14, 0.19, 0.10];

export class ColorPicker {
	private colors: Map<string, string> = new Map();
	constructor(instance: string[]) {
		this.colors = this.setColors(instance);
	}
	setColors(instance: string[]): Map<string, string> {
		let entry = new Map<string, string>();
		const ordered = [...new Set(instance)].sort();
		ordered.forEach((occ: string, index: number) => {
			const hue = (index * GOLDEN_ANGLE) % 360;
			const tier = index % LIGHTNESS.length;
			entry.set(occ, formatHex(oklch({ mode: 'oklch', l: LIGHTNESS[tier]!, c: CHROMA[tier]!, h: hue })));
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
