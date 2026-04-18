<script lang="ts">
	import { page } from '$app/state';
	import { getCollection } from '$lib/api/collections';
	import { auth } from '$lib/auth.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import ShowListItem from '$lib/components/ShowListItem.svelte';
	import type { MovieInfo, ShowInfo, CollectionInfo } from '$lib/types/export_types';

	const collectionId = $derived(page.params.collectionID ?? '');

	let collection = $state<CollectionInfo | null>(null);
	let movies = $state<MovieInfo[]>([]);
	let shows = $state<ShowInfo[]>([]);

	let error_message = $state('');

	async function fetchData() {
		try {
			if (!auth.token) return;
			const result = await getCollection(auth.token, collectionId);
			collection = result.collection;
			movies = result.movies;
			shows = result.shows;
		} catch (e) {
			error_message = (e as Error).message;
		}
	}

	$effect(() => {
		collectionId;
		fetchData();
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
		<h1 class="mb-8 text-2xl font-bold tracking-wide">{collection.name}</h1>
		{#if shows.length > 0}
			<h2 class="mt-8 mb-2 text-xl font-bold tracking-wide">Shows</h2>
			<ItemGrid>
				{#each shows as show (show.id)}
					<ShowListItem {show} />
				{/each}
			</ItemGrid>
		{/if}
		{#if movies.length > 0}
			<h2 class="mt-8 mb-2 text-xl font-bold tracking-wide">Movies</h2>
			<ItemGrid>
				{#each movies as movie (movie.id)}
					<MovieListItem {movie} />
				{/each}
			</ItemGrid>
		{/if}
	</section>
{/if}
