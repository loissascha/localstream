<script lang="ts">
	import { page } from '$app/state';
	import {
		getCollection,
		removeMovieFromCollection,
		removeShowFromCollection
	} from '$lib/api/collections';
	import { auth } from '$lib/auth.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import ShowListItem from '$lib/components/ShowListItem.svelte';
	import type { MovieInfo, ShowInfo, CollectionInfo } from '$lib/types/export_types';

	const collectionId = $derived(page.params.collectionID ?? '');

	let collection = $state<CollectionInfo | null>(null);
	let movies = $state<MovieInfo[]>([]);
	let shows = $state<ShowInfo[]>([]);
	let sortedMovies = $derived(
		[...movies].sort((a, b) => {
			if (a.year < b.year) {
				return -1;
			}
			if (a.year > b.year) {
				return 1;
			}
			return 0;
		})
	);
	let sortedShows = $derived(
		[...shows].sort((a, b) => {
			if (a.year < b.year) {
				return -1;
			}
			if (a.year > b.year) {
				return 1;
			}
			return 0;
		})
	);

	let selectedMovies = $state<Record<string, boolean>>({});
	let selectedShows = $state<Record<string, boolean>>({});
	let selectedCount = $state(0);

	let error_message = $state('');

	async function fetchData() {
		try {
			if (!auth.token) return;
			const result = await getCollection(auth.token, collectionId);
			collection = result.collection;
			movies = result.movies;
			shows = result.shows;
			selectedMovies = Object.fromEntries(movies.map((movie) => [movie.id, false]));
			selectedShows = Object.fromEntries(shows.map((show) => [show.id, false]));
		} catch (e) {
			error_message = (e as Error).message;
		}
	}

	async function removeSelectedMoviesFromCollection() {
		try {
			if (!auth.token) return;
			for (const [id, isSelected] of Object.entries(selectedMovies)) {
				if (isSelected) {
					await removeMovieFromCollection(auth.token, collectionId, id);
				}
			}
			for (const [id, isSelected] of Object.entries(selectedShows)) {
				if (isSelected) {
					await removeShowFromCollection(auth.token, collectionId, id);
				}
			}
			selectedMovies = Object.fromEntries(movies.map((movie) => [movie.id, false]));
			selectedShows = Object.fromEntries(shows.map((show) => [show.id, false]));
			await fetchData();
		} catch (e) {
			alert((e as Error).message);
		}
	}

	$effect(() => {
		collectionId;
		fetchData();
	});

	$effect(() => {
		selectedMovies;
		selectedShows;
		var sc = 0;
		for (const [id, isSelected] of Object.entries(selectedMovies)) {
			if (isSelected) {
				sc++;
			}
		}
		for (const [id, isSelected] of Object.entries(selectedShows)) {
			if (isSelected) {
				sc++;
			}
		}
		selectedCount = sc;
	});
</script>

{#if error_message != ''}
	<div class="my-2 text-red-500">
		{error_message}
	</div>
{/if}

{#if collection == null}
	<div>Loading...</div>
{:else}
	<section class="my-4">
		{#if selectedCount > 0}
			<section
				class="sticky top-0 right-0 left-0 z-10 flex items-center justify-end gap-4 bg-neutral-900 p-4"
			>
				<button
					class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
					onclick={() => {
						removeSelectedMoviesFromCollection();
					}}>Remove all from Collection</button
				>
				<div>{selectedCount} selected</div>
			</section>
		{/if}
		<h1 class="mb-8 text-2xl font-bold tracking-wide">{collection.name}</h1>
		{#if shows.length > 0}
			<h2 class="mt-8 mb-2 text-xl font-bold tracking-wide">Shows</h2>
			<ItemGrid>
				{#each sortedShows as show (show.id)}
					<ShowListItem {show} selectable bind:selected={selectedShows[show.id]} />
				{/each}
			</ItemGrid>
		{/if}
		{#if movies.length > 0}
			<h2 class="mt-8 mb-2 text-xl font-bold tracking-wide">Movies</h2>
			<ItemGrid>
				{#each sortedMovies as movie (movie.id)}
					<MovieListItem
						{movie}
						selectable
						bind:selected={selectedMovies[movie.id]}
						showFinished={movie.finished}
					/>
				{/each}
			</ItemGrid>
		{/if}
	</section>
{/if}
