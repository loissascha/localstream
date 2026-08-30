import { addMovieToCollection } from './api/collections';
import { getMovie, listMovies } from './api/movies';
import {
	deleteWatchstateMovie,
	setWatchstateFinishedMovie,
	updateWatchstateMovie
} from './api/watchstate_movie';
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
	movies.movies = response.movies.sort((a, b) => {
		if (a.name < b.name) return -1;
		if (a.name > b.name) return 1;
		return 0;
	});
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

export async function clearMoviesSelection() {
	movies.selectedMovies = Object.fromEntries(movies.movies.map((movie) => [movie.id, false]));
}

export async function reloadSingleMovie(movieId: string) {
	try {
		if (!auth.token) return;
		const response = await getMovie(auth.token, movieId);
		movies.movies = movies.movies.map((movie) => {
			if (movie.id !== movieId) return movie;
			return response;
		});
	} catch (e) {
		alert((e as Error).message);
	}
}

export async function setSelectedMoviesToWatched() {
	try {
		if (!auth.token) return;
		for (const [id, isSelected] of Object.entries(movies.selectedMovies)) {
			if (isSelected) {
				await setWatchstateFinishedMovie(auth.token, id);
			}
		}
	} catch (e) {
		alert((e as Error).message);
	}
}

export async function deleteWatchstateForSelectedMovies() {
	try {
		if (!auth.token) return;
		for (const [id, isSelected] of Object.entries(movies.selectedMovies)) {
			if (isSelected) {
				await deleteWatchstateMovie(auth.token, id);
			}
		}
	} catch (e) {
		alert((e as Error).message);
	}
}

export async function setMovieWatchstate(movieId: string, position: number, duration: number) {
	try {
		if (!auth.token) return;
		const finished = duration > 0 && position >= Math.max(duration - 10, 0);
		await updateWatchstateMovie(auth.token, {
			movie_id: movieId,
			position: position,
			duration: duration,
			finished: finished
		});

		var percent = 0.0;
		if (duration > 0) {
			percent = (100 / duration) * position;
		}
		if (finished) {
			percent = 100;
		}

		movies.movies = movies.movies.map((movie) => {
			if (movie.id !== movieId) return movie;
			return {
				...movie,
				percentage: percent
			};
		});
	} catch (e) {
		throw e;
	}
}
