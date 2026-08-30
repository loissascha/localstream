<script lang="ts">
	import { resolve } from '$app/paths';
	import type { CollectionInfo } from '$lib/types/export_types';
	import CollectionInfoDisplay from './CollectionInfoDisplay.svelte';
	import ListItemA from './ListItemA.svelte';

	let { collection }: { collection: CollectionInfo } = $props();
</script>

<ListItemA
	href={resolve('/(protected)/(user)/collections/[collectionID]', { collectionID: collection.id })}
>
	<CollectionInfoDisplay {collection} />
</ListItemA>
