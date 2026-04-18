<script lang="ts">
	import { listMovies } from '$lib/api/movies';
	import { auth } from '$lib/auth.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import type { MovieInfo } from '$lib/types/export_types';

	let movies = $state<MovieInfo[]>([]);
	let selectedMovies = $state<Record<string, boolean>>({});
	let loading = $state(true);

	async function loadMovies() {
		try {
			if (!auth.token) return;
			const data = await listMovies(auth.token);
			movies = data.movies;
			movies.sort((a, b) => {
				if (a.name > b.name) return 1;
				if (a.name < b.name) return -1;
				return 0;
			});
			selectedMovies = Object.fromEntries(
				movies.map((movie) => [movie.id, selectedMovies[movie.id] ?? false])
			);
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadMovies();
	});
</script>

<main>
	{#if loading}
		<p>Loading movies...</p>
	{:else}
		<section class="my-8">
			<ItemGrid>
				{#each movies as movie (movie.id)}
					<MovieListItem {movie} selectable bind:selected={selectedMovies[movie.id]} />
				{/each}
			</ItemGrid>
		</section>
	{/if}
</main>
