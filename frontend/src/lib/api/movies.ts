import type { MovieInfo, MovieListResponse } from '$lib/types/export_types';

export async function getMovie(bearerToken: string, movieID: string): Promise<MovieInfo> {
	try {
		const res = await fetch('/api/v1/movies/' + movieID, {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			const m = res.text();
			throw new Error(`Failed to load movie: ${res.status} ${m}`);
		}

		const data = (await res.json()) as MovieInfo;
		return data;
	} catch (error) {
		throw error;
	}
}

export async function listMovies(
	bearerToken: string,
	latest: boolean = false
): Promise<MovieListResponse> {
	try {
		var q = '';
		if (latest) {
			q = '?limit=latest';
		}
		const res = await fetch('/api/v1/movies/list' + q, {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			const m = res.text();
			throw new Error(`Failed to load movies: ${res.status} ${m}`);
		}

		const data = (await res.json()) as MovieListResponse;
		return data;
	} catch (error) {
		throw error;
	}
}
