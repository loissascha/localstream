import type { ShowMetadataInfo } from '$lib/types/export_types';

export async function loadShowMetadata(
	bearerToken: string,
	showId: string
): Promise<ShowMetadataInfo[]> {
	try {
		const res = await fetch('/api/v1/show/metadata/' + showId, {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load metadata: ${res.status}`);
		}

		const data = (await res.json()) as ShowMetadataInfo[];
		return data;
	} catch (error) {
		throw error;
	}
}
