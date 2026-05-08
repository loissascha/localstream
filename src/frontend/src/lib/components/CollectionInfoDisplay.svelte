<script lang="ts">
	import { type CollectionInfo } from '$lib/types/export_types';
	import CollectionPreviewImage from './CollectionPreviewImage.svelte';

	let { collection }: { collection: CollectionInfo } = $props();
</script>

<div
	class="relative flex h-full overflow-hidden p-4 transition-all duration-200 group-hover:border-neutral-500"
>
	<div class="relative flex w-full flex-col gap-4">
		<div class="flex grow flex-col justify-between space-y-2">
			<CollectionPreviewImage {collection} />
			<div class="w-full truncate text-center font-bold text-neutral-200">{collection.name}</div>
		</div>
	</div>
</div>
