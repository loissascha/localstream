import type {
	CreateLibraryRequest,
	CreateLibraryResponse,
	UpdateLibraryRequest,
	UpdateLibraryResponse
} from '$lib/types/export_types';

export async function createLibrary(bearerToken: string, request: CreateLibraryRequest) {
	try {
		const res = await fetch('/api/admin/libraries/create', {
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

export async function updateLibrary(bearerToken: string, request: UpdateLibraryRequest) {
	try {
		const res = await fetch('/api/admin/libraries/update', {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + bearerToken
			},
			body: JSON.stringify(request)
		});
		if (!res.ok) {
			throw new Error(`Failed to update library: ${res.status}`);
		}

		const data = (await res.json()) as UpdateLibraryResponse;
		return data;
	} catch (error) {
		throw error;
	}
}

export async function deleteLibrary(bearerToken: string, id: string) {
	try {
		const res = await fetch('/api/admin/libraries/' + id, {
			method: 'DELETE',
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to delete library: ${res.status}`);
		}
	} catch (error) {
		throw error;
	}
}
