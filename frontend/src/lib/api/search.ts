import type { SearchResponse } from '$lib/types/export_types';

export async function searchLibrary(
	bearerToken: string,
	searchQuery: string
): Promise<SearchResponse> {
	const res = await fetch('/api/v1/search?q=' + encodeURIComponent(searchQuery), {
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to search library: ${res.status} ${m}`);
	}

	const data = (await res.json()) as SearchResponse;
	return data;
}
