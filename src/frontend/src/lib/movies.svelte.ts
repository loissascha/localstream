import { addMovieToCollection } from './api/collections';
import { listMovies } from './api/movies';
import { auth } from './auth.svelte';
import type { MovieInfo } from './types/export_types';

type MoviesState = {
	initialized: boolean;
	movies: MovieInfo[];
	selectedMovies: Record<string, boolean>;
};

export const movies = $state<MoviesState>({
	initialized: false,
	movies: [],
	selectedMovies: {}
});

export async function loadMoviesDatabase() {
	if (!auth.initialized) {
		throw new Error('Auth is not initialized!');
	}
	if (!auth.loggedIn || !auth.token) {
		throw new Error('Not logged in');
	}

	const response = await listMovies(auth.token);
	movies.movies = response.movies;
	movies.selectedMovies = Object.fromEntries(
		movies.movies.map((movie) => [movie.id, movies.selectedMovies[movie.id] ?? false])
	);

	movies.initialized = true;
}

export async function addSelectedMoviesToCollection(collectionId: string) {
	try {
		if (!auth.token) return;
		for (const [id, isSelected] of Object.entries(movies.selectedMovies)) {
			if (isSelected) {
				await addMovieToCollection(auth.token, collectionId, id);
			}
		}
		movies.selectedMovies = Object.fromEntries(movies.movies.map((movie) => [movie.id, false]));
	} catch (e) {
		alert((e as Error).message);
	}
}
