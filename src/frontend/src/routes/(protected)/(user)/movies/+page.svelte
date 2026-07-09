<script lang="ts">
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import SelectCollectionOverlay from '$lib/components/overlays/SelectCollectionOverlay.svelte';
	import {
		addSelectedMoviesToCollection,
		clearMoviesSelection,
		deleteWatchstateForSelectedMovies,
		loadMoviesDatabase,
		movies,
		setSelectedMoviesToWatched
	} from '$lib/movies.svelte';

	const VISIBLE_PER_PAGE = 50;

	let showAddToCollectionOverlay = $state(false);
	let filterByName = $state('');

	let visibleCount = $state(VISIBLE_PER_PAGE);
	let visibleMovies = $derived(
		movies.movies
			.filter((v) => v.name.toLowerCase().includes(filterByName.toLowerCase()))
			.slice(0, visibleCount)
	);
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

	let selectedMoviesCount = $derived(
		Object.entries(movies.selectedMovies).filter(([id, isSelected]) => isSelected).length
	);
</script>

<main>
	<section
		class={`sticky top-0 right-0 left-0 z-50 flex items-center justify-between gap-4 py-4 ${selectedMoviesCount > 0 ? "bg-neutral-900" : ""}`}
	>
		<div>
			<input
				bind:value={filterByName}
				type="text"
				placeholder="Filter by name"
				class="rounded-full border border-transparent bg-neutral-800 px-4 py-2 transition outline-none focus:border-neutral-600"
			/>
		</div>
		{#if selectedMoviesCount > 0}
			<div class="flex flex-wrap items-center justify-end gap-4">
				<button
					class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
					onclick={() => {
						setSelectedMoviesToWatched().then(() => {
							clearMoviesSelection();
							loadMoviesDatabase();
						});
					}}>Mark watched</button
				>
				<button
					class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
					onclick={() => {
						deleteWatchstateForSelectedMovies().then(() => {
							clearMoviesSelection();
							loadMoviesDatabase();
						});
					}}>Mark not watched</button
				>
				<button
					class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
					onclick={() => {
						showAddToCollectionOverlay = true;
					}}>Add to Collection</button
				>
				<button
					class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
					onclick={() => {
						clearMoviesSelection();
					}}>Clear selection</button
				>
				<div>{selectedMoviesCount} selected</div>
			</div>
		{/if}
	</section>
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
