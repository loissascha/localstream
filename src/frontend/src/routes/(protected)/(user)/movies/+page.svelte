<script lang="ts">
	import { addMovieToCollection } from '$lib/api/collections';
	import { listMovies } from '$lib/api/movies';
	import { auth } from '$lib/auth.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import SelectCollectionOverlay from '$lib/components/overlays/SelectCollectionOverlay.svelte';
	import type { MovieInfo } from '$lib/types/export_types';

	let movies = $state<MovieInfo[]>([]);
	let selectedMovies = $state<Record<string, boolean>>({});
	let selectedMoviesCount = $state(0);
	let loading = $state(true);
	let showAddToCollection = $state(false);

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

	async function addSelectedMoviesToCollection(collectionId: string) {
		try {
			if (!auth.token) return;
			for (const [id, isSelected] of Object.entries(selectedMovies)) {
				if (isSelected) {
					await addMovieToCollection(auth.token, collectionId, id);
				}
			}
			selectedMovies = Object.fromEntries(movies.map((movie) => [movie.id, false]));
		} catch (e) {
			alert((e as Error).message);
		} finally {
			showAddToCollection = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadMovies();
	});

	$effect(() => {
		selectedMovies;
		var sc = 0;
		for (const [id, isSelected] of Object.entries(selectedMovies)) {
			if (isSelected) {
				sc++;
			}
		}
		selectedMoviesCount = sc;
	});
</script>

<main>
	{#if loading}
		<p>Loading movies...</p>
	{:else}
		{#if selectedMoviesCount > 0}
			<section
				class="sticky top-0 right-0 left-0 z-10 flex items-center justify-end gap-4 bg-neutral-900 p-4"
			>
				<button
					class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
					onclick={() => {
						showAddToCollection = true;
					}}>Add all to Collection</button
				>
				<div>{selectedMoviesCount} selected</div>
			</section>
		{/if}
		<section class="my-8">
			<ItemGrid>
				{#each movies as movie (movie.id)}
					<MovieListItem
						{movie}
						selectable
						bind:selected={selectedMovies[movie.id]}
						showFinished={movie.finished}
					/>
				{/each}
			</ItemGrid>
		</section>
	{/if}
</main>

{#if showAddToCollection}
	<SelectCollectionOverlay
		selectedCollection={(collectionId) => {
			addSelectedMoviesToCollection(collectionId);
		}}
		close={() => {
			showAddToCollection = false;
		}}
	/>
{/if}
