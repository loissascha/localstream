import { API_URL } from '$lib/consts';
import type { SaveWatchstateRequest, WatchstateResponse } from '$lib/types/export_types';

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
        throw new Error('Error: ' + response.status);
    }
    const result = (await response.json()) as WatchstateResponse;
    return result;
}
