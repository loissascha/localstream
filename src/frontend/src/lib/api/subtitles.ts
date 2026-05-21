import type { SubtitleSupportedLanguage } from '$lib/types/export_types';

export async function loadSupportedSubtitleLanguages(
	bearerToken: string
): Promise<SubtitleSupportedLanguage[]> {
	try {
		const res = await fetch('/api/v1/subtitles/languages', {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load subtitles: ${res.status}`);
		}

		const data = (await res.json()) as SubtitleSupportedLanguage[];
		return data;
	} catch (error) {
		throw error;
	}
}
