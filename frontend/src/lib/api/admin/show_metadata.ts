
export async function setPrimaryMetadataForShow(bearerToken: string, showId: string, metadataId: string) {
	try {
		const res = await fetch('/api/v1/show/metadata/' + showId + '/set/primary/' + metadataId, {
			method: "POST",
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
