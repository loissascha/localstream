<script lang="ts">
	import { type ShowInfo } from '$lib/types/export_types';
	import { auth } from '$lib/auth.svelte';
	import LastWatched from '$lib/components/LastWatched.svelte';
	import LastWatchedMovies from '$lib/components/LastWatchedMovies.svelte';
	import { loadShows } from '$lib/api/shows';
	import ShowListItem from '$lib/components/ShowListItem.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import ShowIcon from '$lib/icons/ShowIcon.svelte';
	import MovieIcon from '$lib/icons/MovieIcon.svelte';
	import { movies } from '$lib/movies.svelte';

	let shows = $state<ShowInfo[]>([]);

	// TODO: need created date in movie info se we can sort properly
	let latestMovies = $derived.by(() => {
		return [...movies.movies]
			.sort((a, b) => {
				if (a.name < b.name) return -1;
				if (a.name > b.name) return 1;
				return 0;
			})
			.slice(0, 10);
	});

	let loadingShows = $state(true);
	let errorMessage = $state('');

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

	<section class="my-8">
		<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
			<MovieIcon />
			Recent Movies
		</h2>
		<ItemGrid>
			{#each latestMovies as movie (movie.id)}
				<MovieListItem {movie} />
			{/each}
		</ItemGrid>
	</section>
</main>
