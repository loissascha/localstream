<script lang="ts">
	import { listCollections } from '$lib/api/collections';
	import { auth } from '$lib/auth.svelte';
	import PlusIcon from '$lib/icons/PlusIcon.svelte';
	import CreateCollectionOverlay from './CreateCollectionOverlay.svelte';
	import Overlay from './Overlay.svelte';

	type CollectionListItem = {
		id: string;
		name: string;
		updated_at: string;
	};

	const formatter = new Intl.DateTimeFormat(undefined, {
		month: 'short',
		year: 'numeric'
	});

	interface Props {
		close: () => void;
		selectedCollection: (collectionId: string) => void;
	}
	let { close, selectedCollection }: Props = $props();

	let showNewCollection = $state(false);

	let collections = $state<CollectionListItem[]>([]);
	const sortedCollections = $derived(
		[...collections].sort((a, b) =>
			a.name.localeCompare(b.name, undefined, { sensitivity: 'base' })
		)
	);
	const collectionCards = $derived(
		sortedCollections.map((collection) => ({
			...collection,
			initial: getInitial(collection.name),
			updatedLabel: formatDate(collection.updated_at)
		}))
	);

	function formatDate(value: string) {
		const parsed = new Date(value);
		if (Number.isNaN(parsed.getTime())) return null;
		return formatter.format(parsed);
	}

	function getInitial(name: string) {
		return name.trim().charAt(0).toUpperCase() || 'C';
	}

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

{#if showNewCollection}
	<CreateCollectionOverlay
		close={() => {
			fetchData().then(() => {
				showNewCollection = false;
			});
		}}
	/>
{:else}
	<Overlay {close}>
		<div class="space-y-2">
			<div class="border-b border-neutral-700/80 pr-10 pb-4">
				<h1 class="text-2xl font-bold tracking-wide text-neutral-100">Collections</h1>
				<p class="mt-1 text-sm text-neutral-400">Choose a collection for this item.</p>
			</div>

			<button
				class="flex gap-2 cursor-pointer"
				onclick={() => {
					showNewCollection = true;
				}}
			>
				<PlusIcon /> New Collection
			</button>

			<div class="space-y-3 pt-2">
				{#if sortedCollections.length === 0}
					<div
						class="rounded-2xl border border-dashed border-neutral-700 bg-neutral-900/60 px-6 py-10 text-center text-neutral-400"
					>
						No collections yet.
					</div>
				{:else}
					{#each collectionCards as collection (collection.id)}
						<button
							type="button"
							onclick={() => {
								selectedCollection(collection.id);
							}}
							class="group flex w-full cursor-pointer items-center gap-4 rounded-2xl border border-neutral-700/80 bg-linear-to-br from-neutral-900 to-neutral-800/90 px-4 py-4 text-left text-neutral-100 shadow-lg shadow-black/10 transition-all duration-200 hover:border-neutral-500 hover:bg-neutral-800 focus-visible:border-brand focus-visible:outline-none"
						>
							<div
								class="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-neutral-100 font-semibold text-neutral-900 shadow-inner shadow-white/20 transition-transform duration-200 group-hover:scale-105"
							>
								{collection.initial}
							</div>

							<div class="min-w-0 flex-1">
								<div class="truncate text-base font-semibold text-neutral-50">
									{collection.name}
								</div>
								{#if collection.updatedLabel}
									<div class="mt-1 text-sm text-neutral-400">Updated {collection.updatedLabel}</div>
								{/if}
							</div>

							<div
								class="text-sm font-medium text-neutral-500 transition-colors duration-200 group-hover:text-neutral-300"
							>
								Select
							</div>
						</button>
					{/each}
				{/if}
			</div>
		</div>
	</Overlay>
{/if}
