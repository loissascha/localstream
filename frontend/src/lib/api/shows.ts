import { API_URL } from '$lib/consts';
import type { ShowListResponse } from '$lib/types/export_types';

export async function loadShows(bearerToken: string): Promise<ShowListResponse> {
	try {
		const res = await fetch(API_URL + '/api/shows', {
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
