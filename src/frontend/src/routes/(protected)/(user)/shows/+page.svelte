<script lang="ts">
	import { addShowToCollection } from '$lib/api/collections';
	import { loadShows } from '$lib/api/shows';
	import { auth } from '$lib/auth.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import SelectCollectionOverlay from '$lib/components/overlays/SelectCollectionOverlay.svelte';
	import ShowListItem from '$lib/components/ShowListItem.svelte';
	import type { ShowInfo } from '$lib/types/export_types';

	let shows = $state<ShowInfo[]>([]);
	let selectedShows = $state<Record<string, boolean>>({});
	let loading = $state(true);
	let errorMessage = $state('');
	let showAddToCollection = $state(false);
	let selectedShowsCount = $state(0);

	async function loadShowsList() {
		try {
			if (!auth.token) return;
			const data = await loadShows(auth.token);
			shows = data.shows;
			shows.sort((a, b) => {
				if (a.name > b.name) return 1;
				if (a.name < b.name) return -1;
				return 0;
			});
			selectedShows = Object.fromEntries(
				shows.map((show) => [show.id, selectedShows[show.id] ?? false])
			);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loading = false;
		}
	}

	async function addSelectedShowsToCollection(collectionId: string) {
		try {
			if (!auth.token) return;
			for (const [id, isSelected] of Object.entries(selectedShows)) {
				if (isSelected) {
					await addShowToCollection(auth.token, collectionId, id);
				}
			}
			selectedShows = Object.fromEntries(shows.map((show) => [show.id, false]));
		} catch (e) {
			alert((e as Error).message);
		} finally {
			showAddToCollection = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadShowsList();
	});

	$effect(() => {
		selectedShows;
		var sc = 0;
		for (const [id, isSelected] of Object.entries(selectedShows)) {
			if (isSelected) {
				sc++;
			}
		}
		selectedShowsCount = sc;
	});
</script>

<main>
	{#if errorMessage}
		<p class="text-red-700">{errorMessage}</p>
	{/if}

	{#if loading}
		<p>Loading shows...</p>
	{:else}
		{#if selectedShowsCount > 0}
			<section class="sticky top-0 right-0 left-0 z-10 flex items-center justify-end gap-4 bg-neutral-900 p-4">
				<button
					class="cursor-pointer rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
					onclick={() => {
						showAddToCollection = true;
					}}>Add all to Collection</button
				>
				<div>{selectedShowsCount} selected</div>
			</section>
		{/if}
		<section class="my-8">
			<ItemGrid>
				{#each shows as show (show.id)}
					<ShowListItem {show} selectable bind:selected={selectedShows[show.id]} />
				{/each}
			</ItemGrid>
		</section>
	{/if}
</main>

{#if showAddToCollection}
	<SelectCollectionOverlay
		selectedCollection={(collectionId) => {
			addSelectedShowsToCollection(collectionId);
		}}
		close={() => {
			showAddToCollection = false;
		}}
	/>
{/if}
