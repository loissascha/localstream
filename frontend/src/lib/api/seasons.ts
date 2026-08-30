import type { SeasonInfo, SeasonListResponse } from '$lib/types/export_types';

export async function loadSeasonsForShow(bearerToken: string, showId: string) {
	try {
		const res = await fetch('/api/v1/seasons/show/' + showId, {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load seasons: ${res.status}`);
		}

		const r = (await res.json()) as SeasonListResponse;
		return r;
	} catch (error) {
		throw error;
	}
}

export async function getSeasonDetails(bearerToken: string, seasonID: string) {
	try {
		const res = await fetch('/api/v1/seasons/' + seasonID, {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load season: ${res.status}`);
		}

		const r = (await res.json()) as SeasonInfo;
		return r;
	} catch (error) {
		throw error;
	}
}
