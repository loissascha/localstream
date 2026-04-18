<script lang="ts">
	import { listCollections } from '$lib/api/collections';
	import { auth } from '$lib/auth.svelte';
	import CollectionListItem from '$lib/components/CollectionListItem.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import CreateCollectionOverlay from '$lib/components/overlays/CreateCollectionOverlay.svelte';
	import PlusIcon from '$lib/icons/PlusIcon.svelte';
	import type { CollectionInfo } from '$lib/types/export_types';

	let collections = $state<CollectionInfo[]>([]);
	let createCollectionOverlayOpen = $state(false);

	async function loadData() {
		if (!auth.token) return;
		const result = await listCollections(auth.token);
		collections = result.collections;
	}

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.token) return;
		loadData();
	});
</script>

<section class="mb-4 flex items-center justify-center gap-4">
	<div>
		{collections.length} Collections
	</div>
	<button
		onclick={() => {
			createCollectionOverlayOpen = true;
		}}
		class="flex cursor-pointer gap-1 rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
	>
		<PlusIcon />
		Create Collection
	</button>
</section>

<section>
	<ItemGrid>
		{#each collections as collection (collection.id)}
			<CollectionListItem {collection} />
		{/each}
	</ItemGrid>
</section>

{#if createCollectionOverlayOpen}
	<CreateCollectionOverlay
		close={() => {
			createCollectionOverlayOpen = false;
			loadData();
		}}
	/>
{/if}
