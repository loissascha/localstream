<script lang="ts">
	import { resolve } from '$app/paths';
	import {
		type ShowInfo,
		type LibraryListItem,
		type MovieInfo
	} from '$lib/types/export_types';
	import { auth } from '$lib/auth.svelte';
	import LastWatched from '$lib/components/LastWatched.svelte';
	import { loadLibraries } from '$lib/api/libraries';
	import LibraryIcon from '$lib/icons/LibraryIcon.svelte';
	import { listMovies } from '$lib/api/movies';
	import LastWatchedMovies from '$lib/components/LastWatchedMovies.svelte';
	import { loadShows } from '$lib/api/shows';
	import ShowListItem from '$lib/components/ShowListItem.svelte';

	let libraries = $state<LibraryListItem[]>([]);
	let shows = $state<ShowInfo[]>([]);
	let movies = $state<MovieInfo[]>([]);
	let selectedLibrary = $state<LibraryListItem | null>(null);
	let loading = $state(true);
	let loadingShows = $state(true);
	let loadingMovies = $state(true);
	let errorMessage = $state('');

	async function loadMovies() {
		try {
			if (!auth.token) return;
			const data = await listMovies(auth.token);
			movies = data.movies;
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		} finally {
			loadingMovies = false;
		}
	}

	async function loadShowsList() {
		try {
			if (!auth.token) return;
			const data = await loadShows(auth.token);
			shows = data.shows;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loadingShows = false;
		}
	}

	async function loadLibrariesData() {
		try {
			if (!auth.token) {
				throw new Error('Auth token not loaded.');
			}
			const data = await loadLibraries(auth.token);
			libraries = data.libraries;
			console.log('libraries', data);

			if (selectedLibrary == null) {
				if (libraries.length > 0) {
					selectLibrary(libraries[0]);
				}
			}
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loading = false;
		}
	}

	function selectLibrary(lib: LibraryListItem | null) {
		selectedLibrary = lib;
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadLibrariesData();
		loadShowsList();
		loadMovies();
	});
</script>

<main>
	{#if errorMessage}
		<p class="text-red-700">{errorMessage}</p>
	{/if}

	<LastWatched />
	<LastWatchedMovies />

	{#if loadingShows}
		<p>Loading shows...</p>
	{:else}
		<section class="my-8">
			<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
				<LibraryIcon />
				Shows
			</h2>
			<div class="flex gap-4">
				{#each shows as show (show.id)}
					<ShowListItem show={show} />
				{/each}
			</div>
		</section>
	{/if}

	{#if loadingMovies}
		<p>Loading movies...</p>
	{:else}
		<section class="my-8">
			<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
				<LibraryIcon />
				Movies
			</h2>
			<div class="flex gap-4">
				{#each movies as movie (movie.id)}
					<a
						href={resolve('/(protected)/(watch)/movies/[movieID]', { movieID: movie.id })}
						class="h-20 w-60 cursor-pointer truncate rounded-lg bg-neutral-800 p-4 hover:bg-neutral-700"
					>
						{movie.name}
					</a>
				{/each}
			</div>
		</section>
	{/if}
</main>
