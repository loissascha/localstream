import { API_URL } from '$lib/consts';
import type {
    SaveWatchstateRequest,
    WatchstateListResponse,
    WatchstateMovieResponse,
    WatchstateMoviesListResponse,
    WatchstateResponse
} from '$lib/types/export_types';

export async function setWatchstateFinished(
    bearerToken: string,
    episodeID: string
): Promise<WatchstateResponse> {
    const response = await fetch(API_URL + '/api/watchstate/episode/' + episodeID + '/finished', {
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
    const response = await fetch(API_URL + '/api/watchstate', {
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

export async function listLatestWatchstateByMovie(bearerToken: string): Promise<WatchstateMovieResponse[]> {
    const response = await fetch(API_URL + '/api/v1/watchstate/movie/latest', {
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

export async function listLatestWatchstateByShow(
    bearerToken: string
): Promise<WatchstateResponse[]> {
    const response = await fetch(API_URL + '/api/watchstate/latest/shows', {
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
    const response = await fetch(API_URL + '/api/watchstate/episode/' + episodeId, {
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
