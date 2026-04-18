import type {
	CollectionDetailResponse,
	CollectionInfo,
	CollectionListResponse,
	CreateCollectionRequest,
	UpdateCollectionRequest
} from '$lib/types/export_types';

export async function listCollections(bearerToken: string): Promise<CollectionListResponse> {
	const res = await fetch('/api/v1/collections', {
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to load collections: ${res.status} ${m}`);
	}

	const data = (await res.json()) as CollectionListResponse;
	return data;
}

export async function getCollection(
	bearerToken: string,
	collectionID: string
): Promise<CollectionDetailResponse> {
	const res = await fetch('/api/v1/collections/' + collectionID, {
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to load collection: ${res.status} ${m}`);
	}

	const data = (await res.json()) as CollectionDetailResponse;
	return data;
}

export async function createCollection(
	bearerToken: string,
	body: CreateCollectionRequest
): Promise<CollectionInfo> {
	const res = await fetch('/api/v1/collections', {
		method: 'POST',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		},
		body: JSON.stringify(body)
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to create collection: ${res.status} ${m}`);
	}

	const data = (await res.json()) as CollectionInfo;
	return data;
}

export async function updateCollection(
	bearerToken: string,
	collectionID: string,
	body: UpdateCollectionRequest
): Promise<CollectionInfo> {
	const res = await fetch('/api/v1/collections/' + collectionID + '/update', {
		method: 'POST',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		},
		body: JSON.stringify(body)
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to update collection: ${res.status} ${m}`);
	}

	const data = (await res.json()) as CollectionInfo;
	return data;
}

export async function deleteCollection(bearerToken: string, collectionID: string): Promise<void> {
	const res = await fetch('/api/v1/collections/' + collectionID, {
		method: 'DELETE',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to delete collection: ${res.status} ${m}`);
	}
}

export async function addMovieToCollection(
	bearerToken: string,
	collectionID: string,
	movieID: string
): Promise<void> {
	const res = await fetch('/api/v1/collections/' + collectionID + '/movies/' + movieID, {
		method: 'POST',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to add movie to collection: ${res.status} ${m}`);
	}
}

export async function removeMovieFromCollection(
	bearerToken: string,
	collectionID: string,
	movieID: string
): Promise<void> {
	const res = await fetch('/api/v1/collections/' + collectionID + '/movies/' + movieID, {
		method: 'DELETE',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to remove movie from collection: ${res.status} ${m}`);
	}
}

export async function addShowToCollection(
	bearerToken: string,
	collectionID: string,
	showID: string
): Promise<void> {
	const res = await fetch('/api/v1/collections/' + collectionID + '/shows/' + showID, {
		method: 'POST',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to add show to collection: ${res.status} ${m}`);
	}
}

export async function removeShowFromCollection(
	bearerToken: string,
	collectionID: string,
	showID: string
): Promise<void> {
	const res = await fetch('/api/v1/collections/' + collectionID + '/shows/' + showID, {
		method: 'DELETE',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (!res.ok) {
		const m = await res.text();
		throw new Error(`Failed to remove show from collection: ${res.status} ${m}`);
	}
}
