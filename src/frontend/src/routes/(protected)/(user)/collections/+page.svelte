<script lang="ts">
	import { listCollections } from '$lib/api/collections';
	import { auth } from '$lib/auth.svelte';
	import CreateCollectionOverlay from '$lib/components/overlays/CreateCollectionOverlay.svelte';
	import PlusIcon from '$lib/icons/PlusIcon.svelte';
	import type { CollectionInfo } from '$lib/types/export_types';

	let collections = $state<CollectionInfo[]>([]);

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

<section class="flex items-center justify-center gap-4">
	<div>
		{collections.length} Collections
	</div>
	<button
		class="flex cursor-pointer gap-1 rounded-full bg-neutral-800 px-4 py-2 hover:bg-neutral-700"
	>
		<PlusIcon />
		Create Collection
	</button>
</section>

<CreateCollectionOverlay
	close={() => {
		alert('close');
	}}
/>
