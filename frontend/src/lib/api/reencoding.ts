import type { ReencodeResponse, ReencodeStatus } from '$lib/types/export_types';

export async function reencodeMovie(bearerToken: string, movieID: string): Promise<ReencodeStatus> {
	try {
		const res = await fetch('/api/v1/reencoding/movie/' + movieID, {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + bearerToken,
				'Content-Type': 'application/json'
			}
		});
		if (!res.ok) {
			const m = res.text();
			throw new Error(`Failed to load data: ${res.status} ${m}`);
		}

		const data = (await res.json()) as ReencodeStatus;
		return data;
	} catch (error) {
		throw error;
	}
}

export async function listReencodingMovies(bearerToken: string): Promise<ReencodeResponse> {
	try {
		const res = await fetch('/api/v1/reencoding/movies', {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			const m = res.text();
			throw new Error(`Failed to load data: ${res.status} ${m}`);
		}

		const data = (await res.json()) as ReencodeResponse;
		return data;
	} catch (error) {
		throw error;
	}
}
