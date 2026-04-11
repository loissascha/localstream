<script lang="ts">
	import { listMovies } from '$lib/api/movies';
	import { auth } from '$lib/auth.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import LibraryIcon from '$lib/icons/LibraryIcon.svelte';
	import type { MovieInfo } from '$lib/types/export_types';

	let movies = $state<MovieInfo[]>([]);
	let loading = $state(true);

	async function loadMovies() {
		try {
			if (!auth.token) return;
			const data = await listMovies(auth.token);
			movies = data.movies;
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
			<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
				<LibraryIcon />
				Movies
			</h2>
			<ItemGrid>
				{#each movies as movie (movie.id)}
					<MovieListItem {movie} />
				{/each}
			</ItemGrid>
		</section>
	{/if}
</main>
