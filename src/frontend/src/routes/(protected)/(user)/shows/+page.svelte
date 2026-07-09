<script lang="ts">
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import SelectCollectionOverlay from '$lib/components/overlays/SelectCollectionOverlay.svelte';
	import ShowListItem from '$lib/components/ShowListItem.svelte';
	import { addSelectedShowsToCollection, clearShowsSelection, shows } from '$lib/shows.svelte';

	const VISIBLE_PER_PAGE = 50;

	let showAddToCollection = $state(false);
	let filterByName = $state('');

	let visibleCount = $state(VISIBLE_PER_PAGE);
	let visibleShows = $derived(
		shows.shows
			.filter((v) => v.name.toLowerCase().includes(filterByName.toLowerCase()))
			.slice(0, visibleCount)
	);
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

	let selectedShowsCount = $derived(
		Object.entries(shows.selectedShows).filter(([id, isSelected]) => isSelected).length
	);
</script>

<main>
	<section
		class={`sticky top-0 right-0 left-0 z-50 flex items-center justify-between gap-4 py-4 ${selectedShowsCount > 0 ? "bg-neutral-900" : ""}`}
	>
		<div>
			<input
				bind:value={filterByName}
				type="text"
				placeholder="Filter by name"
				class="rounded-full border border-transparent bg-neutral-800 px-4 py-2 transition outline-none focus:border-neutral-600"
			/>
		</div>
		{#if selectedShowsCount > 0}
			<div class="flex flex-wrap items-center justify-end gap-4">
				<button
					class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
					onclick={() => {
						showAddToCollection = true;
					}}>Add to Collection</button
				>
				<button
					class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
					onclick={() => {
						clearShowsSelection();
					}}>Clear selection</button
				>
				<div>{selectedShowsCount} selected</div>
			</div>
		{/if}
	</section>
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
