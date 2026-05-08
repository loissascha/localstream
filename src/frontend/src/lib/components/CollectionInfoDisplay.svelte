<script lang="ts">
	import { type CollectionInfo } from '$lib/types/export_types';
	import CollectionPreviewImage from './CollectionPreviewImage.svelte';

	let { collection }: { collection: CollectionInfo } = $props();

	const formatter = new Intl.DateTimeFormat(undefined, {
		month: 'short',
		year: 'numeric'
	});

	function formatDate(value: string) {
		const parsed = new Date(value);
		if (Number.isNaN(parsed.getTime())) return null;
		return formatter.format(parsed);
	}

	const updatedLabel = $derived(formatDate(collection.updated_at));
</script>

<div
	class="via-neutral-850 relative flex h-full overflow-hidden rounded-xl border border-neutral-700/80 bg-linear-to-br from-neutral-900 to-neutral-900 p-4 text-neutral-200 shadow-lg shadow-black/20 transition-colors duration-200 group-hover:border-neutral-500"
>
	<div class="relative flex w-full flex-col gap-4">
		<div class="flex grow flex-col justify-between space-y-2">
			<div class="flex items-start gap-3">
				<div class="min-w-0">
					<div class="line-clamp-2 text-lg leading-tight font-bold text-neutral-100">
						{collection.name}
					</div>
				</div>
			</div>
			<CollectionPreviewImage {collection} />
			{#if updatedLabel}
				<div class="mt-1 text-right text-sm text-neutral-400">Updated {updatedLabel}</div>
			{/if}
		</div>
	</div>
</div>
