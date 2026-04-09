import { API_URL } from "$lib/consts";
import type { LibraryListResponse } from "$lib/types/export_types";

export async function loadLibraries(bearerToken: string) {
	try {
		const res = await fetch(API_URL + '/api/libraries', {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load videos: ${res.status}`);
		}

		const data = (await res.json()) as LibraryListResponse;
		return data
	} catch (error) {
		throw error;
	}
}
