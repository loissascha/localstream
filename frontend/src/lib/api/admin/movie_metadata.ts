import { API_URL } from "$lib/consts";

export async function setPrimaryMetadataForMovie(
	bearerToken: string,
	movieId: string,
	metadataId: string
) {
	try {
		const res = await fetch(API_URL + '/api/v1/movie/metadata/' + movieId + '/set/primary/' + metadataId, {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to set primary metadata: ${res.status}`);
		}
	} catch (error) {
		throw error;
	}
}
