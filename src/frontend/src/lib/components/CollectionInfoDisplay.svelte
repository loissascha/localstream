<script lang="ts">
	import { type CollectionInfo } from '$lib/types/export_types';

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
	const initial = $derived(collection.name.trim().charAt(0).toUpperCase() || 'C');
</script>

<div
	class="via-neutral-850 relative flex h-40 overflow-hidden rounded-xl border border-neutral-700/80 bg-linear-to-br from-neutral-900 to-neutral-900 p-4 text-neutral-200 shadow-lg shadow-black/20 transition-colors duration-200 group-hover:border-neutral-500"
>
	<div class="relative flex w-full flex-col gap-4">
		<div class="flex items-start justify-between gap-3">
			<div
				class="rounded-full border border-neutral-600/80 bg-neutral-800/80 px-3 py-1 text-xs font-semibold tracking-[0.24em] text-neutral-400 uppercase"
			>
				Collection
			</div>
		</div>

		<div class="flex grow flex-col justify-between space-y-2">
			<div class="flex items-start gap-3">
				<div
					class="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl bg-neutral-100 font-semibold text-neutral-900 shadow-inner shadow-white/20"
				>
					{initial}
				</div>

				<div class="min-w-0">
					<div class="line-clamp-2 text-lg leading-tight font-bold text-neutral-100">
						{collection.name}
					</div>
				</div>
			</div>
			{#if updatedLabel}
				<div class="mt-1 text-right text-sm text-neutral-400">Updated {updatedLabel}</div>
			{/if}
		</div>
	</div>
</div>
