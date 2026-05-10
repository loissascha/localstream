import type {
	SaveMovieWatchstateRequest,
	WatchstateInfo,
	WatchstateMovieResponse,
	WatchstateMoviesListResponse
} from '$lib/types/export_types';

export async function updateWatchstateMovie(
	bearerToken: string,
	body: SaveMovieWatchstateRequest
): Promise<WatchstateInfo> {
	const response = await fetch('/api/v1/movie/watchstate', {
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
	const result = (await response.json()) as WatchstateInfo;
	return result;
}

export async function setWatchstateFinishedMovie(
	bearerToken: string,
	movieID: string
): Promise<WatchstateMovieResponse> {
	const response = await fetch('/api/v1/watchstate/movie/' + movieID + '/finished', {
		method: 'POST',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
	const result = (await response.json()) as WatchstateMovieResponse;
	return result;
}

export async function deleteWatchstateMovie(bearerToken: string, movieId: string) {
	const response = await fetch('/api/v1/watchstate/movie/' + movieId + '/delete', {
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

export async function listLatestWatchstateByMovie(
	bearerToken: string
): Promise<WatchstateMovieResponse[]> {
	const response = await fetch('/api/v1/watchstate/movie/latest', {
		method: 'GET',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
	const result = (await response.json()) as WatchstateMoviesListResponse;
	return result.watchstates;
}

export async function getWatchstateForMovie(
	bearerToken: string,
	movieId: string
): Promise<WatchstateInfo> {
	const response = await fetch('/api/v1/watchstate/movie/' + movieId, {
		method: 'GET',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('' + response.status);
	}
	const result = (await response.json()) as WatchstateInfo;
	return result;
}
