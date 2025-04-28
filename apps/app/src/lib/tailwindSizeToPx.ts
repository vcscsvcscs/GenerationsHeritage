export function tailwindClassToPixels(className: string): number | null {
	const remSize = getRemInPixels(); // <-- real rem size at runtime

	const regex = /^(w|h)-(\d+)$/;
	const match = className.match(regex);
	if (!match) return null;

	const value = parseInt(match[2], 10);
	return (value / 4) * remSize;
}

export function getRemInPixels(): number {
	try {
		const fontSize = getComputedStyle(document.documentElement).fontSize;
		return parseFloat(fontSize);
	} catch (e) {
		return 16; // Default to 16px if unable to get computed style
	}
}
