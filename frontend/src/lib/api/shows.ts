import type { ShowListResponse } from '$lib/types/export_types';

export async function loadShows(
	bearerToken: string,
	latest: boolean = false
): Promise<ShowListResponse> {
	try {
		var q = '';
		if (latest) {
			q = '?limit=latest';
		}
		const res = await fetch('/api/shows' + q, {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load shows: ${res.status}`);
		}

		const data = (await res.json()) as ShowListResponse;
		return data;
	} catch (error) {
		throw error;
	}
}
