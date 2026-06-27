import { addShowToCollection } from './api/collections';
import { loadShows } from './api/shows';
import { updateWatchstate } from './api/watchstate';
import { auth } from './auth.svelte';
import type { ShowInfo } from './types/export_types';

type ShowsState = {
	initialized: boolean;
	shows: ShowInfo[];
	selectedShows: Record<string, boolean>;
};

export const shows = $state<ShowsState>({
	initialized: false,
	shows: [],
	selectedShows: {}
});

export async function loadShowsDatabase() {
	if (!auth.initialized) {
		throw new Error('Auth is not initialized!');
	}
	if (!auth.loggedIn || !auth.token) {
		throw new Error('Not logged in');
	}

	const response = await loadShows(auth.token);
	shows.shows = response.shows.sort((a, b) => {
		if (a.name < b.name) return -1;
		if (a.name > b.name) return 1;
		return 0;
	});
	shows.selectedShows = Object.fromEntries(
		shows.shows.map((show) => [show.id, shows.selectedShows[show.id] ?? false])
	);

	shows.initialized = true;
}

export async function clearShowsSelection() {
	shows.selectedShows = Object.fromEntries(shows.shows.map((show) => [show.id, false]));
}

export async function addSelectedShowsToCollection(collectionId: string) {
	try {
		if (!auth.token) return;
		for (const [id, isSelected] of Object.entries(shows.selectedShows)) {
			if (isSelected) {
				await addShowToCollection(auth.token, collectionId, id);
			}
		}
		shows.selectedShows = Object.fromEntries(shows.shows.map((show) => [show.id, false]));
	} catch (e) {
		alert((e as Error).message);
	}
}

export async function setShowWatchstate(
	showId: string,
	seasonId: string,
	episodeId: string,
	position: number,
	duration: number
) {
	try {
		if (!auth.token) return;
		const finished = duration > 0 && position >= Math.max(duration - 10, 0);
		const normalizedDuration = Number.isFinite(duration) ? Number(duration.toFixed(2)) : 0;
		await updateWatchstate(auth.token, {
			episode_id: episodeId,
			season_id: seasonId,
			show_id: showId,
			position: position,
			duration: normalizedDuration,
			finished: finished
		});

		var percent = 0.0;
		if (duration > 0) {
			percent = (100 / duration) * position;
		}
		if (finished) {
			percent = 100;
		}

		shows.shows = shows.shows.map((show) => {
			if (show.id !== showId) return show;
			return {
				...show,
				percentage: percent
			};
		});
	} catch (e) {
		throw e;
	}
}
