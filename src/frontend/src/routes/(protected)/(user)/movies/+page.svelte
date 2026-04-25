<script lang="ts">
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import SelectCollectionOverlay from '$lib/components/overlays/SelectCollectionOverlay.svelte';
	import { addSelectedMoviesToCollection, movies } from '$lib/movies.svelte';

	let selectedMoviesCount = $state(0);
	let showAddToCollection = $state(false);

	$effect(() => {
		movies.selectedMovies;
		var sc = 0;
		for (const [id, isSelected] of Object.entries(movies.selectedMovies)) {
			if (isSelected) {
				sc++;
			}
		}
		selectedMoviesCount = sc;
	});
</script>

<main>
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
			{#each movies.movies as movie (movie.id)}
				<MovieListItem {movie} selectable bind:selected={movies.selectedMovies[movie.id]} />
			{/each}
		</ItemGrid>
	</section>
</main>

{#if showAddToCollection}
	<SelectCollectionOverlay
		selectedCollection={(collectionId) => {
			addSelectedMoviesToCollection(collectionId)
				.then(() => {
					showAddToCollection = false;
				})
				.catch((e) => {
					const m = (e as Error).message;
				});
		}}
		close={() => {
			showAddToCollection = false;
		}}
	/>
{/if}
