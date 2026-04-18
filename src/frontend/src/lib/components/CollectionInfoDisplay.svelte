<script lang="ts">
	import { type CollectionInfo } from '$lib/types/export_types';

	let { collection }: { collection: CollectionInfo } = $props();
</script>

<div class="h-40 rounded border border-neutral-500 p-4">
	<button class={`my-1  w-full cursor-pointer font-bold text-neutral-200`}>
		<div class="w-full text-center">{collection.name}</div>
	</button>
</div>
