import type {
	SaveWatchstateRequest,
	WatchstateListResponse,
	WatchstateResponse
} from '$lib/types/export_types';

export async function deleteWatchstate(bearerToken: string, episodeID: string) {
	const response = await fetch('/api/watchstate/episode/' + episodeID + '/delete', {
		method: 'DELETE',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
}

export async function setWatchstateFinished(
	bearerToken: string,
	episodeID: string
): Promise<WatchstateResponse> {
	const response = await fetch('/api/watchstate/episode/' + episodeID + '/finished', {
		method: 'POST',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
	const result = (await response.json()) as WatchstateResponse;
	return result;
}

export async function updateWatchstate(
	bearerToken: string,
	body: SaveWatchstateRequest
): Promise<WatchstateResponse> {
	const response = await fetch('/api/watchstate', {
		method: 'POST',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		},
		body: JSON.stringify(body)
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
	const result = (await response.json()) as WatchstateResponse;
	return result;
}

export async function listLatestWatchstateByShow(
	bearerToken: string
): Promise<WatchstateResponse[]> {
	const response = await fetch('/api/watchstate/latest/shows', {
		method: 'GET',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
	const result = (await response.json()) as WatchstateListResponse;
	return result.watchstates;
}

export async function getWatchstateForEpisode(
	bearerToken: string,
	episodeId: string
): Promise<WatchstateResponse> {
	const response = await fetch('/api/watchstate/episode/' + episodeId, {
		method: 'GET',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('' + response.status);
	}
	const result = (await response.json()) as WatchstateResponse;
	return result;
}
