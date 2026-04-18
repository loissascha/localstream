<script lang="ts">
	import { type ShowInfo, type MovieInfo } from '$lib/types/export_types';
	import { auth } from '$lib/auth.svelte';
	import LastWatched from '$lib/components/LastWatched.svelte';
	import { listMovies } from '$lib/api/movies';
	import LastWatchedMovies from '$lib/components/LastWatchedMovies.svelte';
	import { loadShows } from '$lib/api/shows';
	import ShowListItem from '$lib/components/ShowListItem.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import ShowIcon from '$lib/icons/ShowIcon.svelte';
	import MovieIcon from '$lib/icons/MovieIcon.svelte';

	let shows = $state<ShowInfo[]>([]);
	let movies = $state<MovieInfo[]>([]);
	let loadingShows = $state(true);
	let loadingMovies = $state(true);
	let errorMessage = $state('');

	async function loadMovies() {
		try {
			if (!auth.token) return;
			const data = await listMovies(auth.token, true);
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
			const data = await loadShows(auth.token, true);
			shows = data.shows;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loadingShows = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
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
				<ShowIcon />
				Recent Shows
			</h2>
			<ItemGrid>
				{#each shows as show (show.id)}
					<ShowListItem {show} />
				{/each}
			</ItemGrid>
		</section>
	{/if}

	{#if loadingMovies}
		<p>Loading movies...</p>
	{:else}
		<section class="my-8">
			<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
				<MovieIcon />
				Recent Movies
			</h2>
			<ItemGrid>
				{#each movies as movie (movie.id)}
					<MovieListItem {movie} showFinished={movie.finished} />
				{/each}
			</ItemGrid>
		</section>
	{/if}
</main>
