<script lang="ts">
	import { listCollections } from '$lib/api/collections';
	import { auth } from '$lib/auth.svelte';
	import type { CollectionInfo } from '$lib/types/export_types';
	import Overlay from './Overlay.svelte';

	interface Props {
		close: () => void;
		selectedCollection: (collectionId: string) => void;
	}
	let { close, selectedCollection }: Props = $props();

	let collections = $state<CollectionInfo[]>([]);

	async function fetchData() {
		try {
			if (!auth.token) return;
			const result = await listCollections(auth.token);
			collections = result.collections;
		} catch (e) {
			alert((e as Error).message);
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.token) return;
		fetchData();
	});
</script>

<Overlay {close}>
	<h1 class="text-2xl font-bold tracking-wide">Collections</h1>
	<div class="my-4">
		{#each collections as collection (collection.id)}
			<button
				onclick={() => {
					selectedCollection(collection.id);
				}}
				class="my-2 block cursor-pointer hover:scale-105 hover:font-bold"
			>
				{collection.name}
			</button>
		{/each}
	</div>
</Overlay>
