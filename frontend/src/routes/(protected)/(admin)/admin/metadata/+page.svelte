<script lang="ts">
	import { listMovies } from '$lib/api/movies';
	import { loadShows } from '$lib/api/shows';
	import { auth } from '$lib/auth.svelte';
	import AdminMovieMetadataBlock from '$lib/components/admin/AdminMovieMetadataBlock.svelte';
	import AdminShowMetadataBlock from '$lib/components/admin/AdminShowMetadataBlock.svelte';
	import { type MovieInfo, type ShowInfo } from '$lib/types/export_types';

	var shows = $state<ShowInfo[]>([]);
	var movies = $state<MovieInfo[]>([]);
	var loadingShows = $state(true);
	var loadingMovies = $state(true);
	var errorMessage = $state('');

	async function loadShowsList() {
		try {
			if (!auth.token) return;
			const data = await loadShows(auth.token);
			shows = data.shows;
		} catch (e) {
			errorMessage = (e as Error).message;
		} finally {
			loadingShows = false;
		}
	}

	async function loadMoviesList() {
		try {
			if (!auth.token) return;
			const data = await listMovies(auth.token);
			movies = data.movies;
		} catch (e) {
			errorMessage = (e as Error).message;
		} finally {
			loadingMovies = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadShowsList();
		loadMoviesList();
	});
</script>

{#if errorMessage != ''}
	<p class="text-red-500">{errorMessage}</p>
{/if}

<section class="grid grid-cols-3 gap-4">
	{#each shows as show}
		<AdminShowMetadataBlock {show} />
	{/each}
</section>
<section class="mt-4 grid grid-cols-3 gap-4">
	{#each movies as movie}
		<AdminMovieMetadataBlock {movie} />
	{/each}
</section>
