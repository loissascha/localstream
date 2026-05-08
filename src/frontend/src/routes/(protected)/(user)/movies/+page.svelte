<script lang="ts">
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import SelectCollectionOverlay from '$lib/components/overlays/SelectCollectionOverlay.svelte';
	import { addSelectedMoviesToCollection, clearMoviesSelection, movies } from '$lib/movies.svelte';

	const VISIBLE_PER_PAGE = 50;

	let showAddToCollectionOverlay = $state(false);

	let visibleCount = $state(VISIBLE_PER_PAGE);
	let visibleMovies = $derived(movies.movies.slice(0, visibleCount));
	let sentinel = $state<HTMLDivElement | null>(null);

	function loadMore() {
		if (visibleCount < movies.movies.length) {
			visibleCount += VISIBLE_PER_PAGE;
		}
	}

	$effect(() => {
		if (!sentinel) return;

		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0]?.isIntersecting) {
					loadMore();
				}
			},
			{
				rootMargin: '600px'
			}
		);

		observer.observe(sentinel);

		return () => observer.disconnect();
	});

	let selectedMoviesCount = $state(0);
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
			class="sticky top-0 right-0 left-0 z-50 flex items-center justify-end gap-4 bg-neutral-900 p-4"
		>
			<button
				class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
				onclick={() => {
					clearMoviesSelection();
				}}>Clear selection</button
			>
			<button
				class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
				onclick={() => {
					showAddToCollectionOverlay = true;
				}}>Add to Collection</button
			>
			<div>{selectedMoviesCount} selected</div>
		</section>
	{/if}
	<section class="my-8">
		<ItemGrid>
			{#each visibleMovies as movie (movie.id)}
				<MovieListItem {movie} selectable bind:selected={movies.selectedMovies[movie.id]} />
			{/each}
		</ItemGrid>
		<div bind:this={sentinel}></div>
	</section>
</main>

{#if showAddToCollectionOverlay}
	<SelectCollectionOverlay
		selectedCollection={(collectionId) => {
			addSelectedMoviesToCollection(collectionId)
				.then(() => {
					showAddToCollectionOverlay = false;
				})
				.catch((e) => {
					const m = (e as Error).message;
				});
		}}
		close={() => {
			showAddToCollectionOverlay = false;
		}}
	/>
{/if}
