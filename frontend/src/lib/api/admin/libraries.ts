import { API_URL } from '$lib/consts';
import type { CreateLibraryRequest, CreateLibraryResponse } from '$lib/types/export_types';

export async function createLibrary(bearerToken: string, request: CreateLibraryRequest) {
	try {
		const res = await fetch(API_URL + '/api/admin/libraries/create', {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + bearerToken
			},
			body: JSON.stringify(request)
		});
		if (!res.ok) {
			throw new Error(`Failed to create library: ${res.status}`);
		}

		const data = (await res.json()) as CreateLibraryResponse;
		return data;
	} catch (error) {
		throw error;
	}
}
