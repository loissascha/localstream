<script lang="ts">
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import SelectCollectionOverlay from '$lib/components/overlays/SelectCollectionOverlay.svelte';
	import ShowListItem from '$lib/components/ShowListItem.svelte';
	import { addSelectedShowsToCollection, clearShowsSelection, shows } from '$lib/shows.svelte';

	const VISIBLE_PER_PAGE = 50;

	let showAddToCollection = $state(false);
	let selectedShowsCount = $state(0);

	let visibleCount = $state(VISIBLE_PER_PAGE);
	let visibleShows = $derived(shows.shows.slice(0, visibleCount));
	let sentinel = $state<HTMLDivElement | null>(null);

	function loadMore() {
		if (visibleCount < shows.shows.length) {
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

	$effect(() => {
		shows.selectedShows;
		var sc = 0;
		for (const [id, isSelected] of Object.entries(shows.selectedShows)) {
			if (isSelected) {
				sc++;
			}
		}
		selectedShowsCount = sc;
	});
</script>

<main>
	{#if selectedShowsCount > 0}
		<section
			class="sticky top-0 right-0 left-0 z-50 flex items-center justify-end gap-4 bg-neutral-900 p-4"
		>
			<button
				class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
				onclick={() => {
					clearShowsSelection();
				}}>Clear selection</button
			>
			<button
				class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
				onclick={() => {
					showAddToCollection = true;
				}}>Add to Collection</button
			>
			<div>{selectedShowsCount} selected</div>
		</section>
	{/if}
	<section class="my-8">
		<ItemGrid>
			{#each visibleShows as show (show.id)}
				<ShowListItem {show} selectable bind:selected={shows.selectedShows[show.id]} />
			{/each}
		</ItemGrid>
		<div bind:this={sentinel}></div>
	</section>
</main>

{#if showAddToCollection}
	<SelectCollectionOverlay
		selectedCollection={(collectionId) => {
			addSelectedShowsToCollection(collectionId).then(() => {
				showAddToCollection = false;
			});
		}}
		close={() => {
			showAddToCollection = false;
		}}
	/>
{/if}
