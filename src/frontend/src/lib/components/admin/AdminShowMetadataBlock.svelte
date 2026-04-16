<script lang="ts">
	import { setPrimaryMetadataForShow } from '$lib/api/admin/show_metadata';
	import { loadShowMetadata } from '$lib/api/show_metadata';
	import { auth } from '$lib/auth.svelte';
	import { type ShowMetadataInfo, type ShowInfo } from '$lib/types/export_types';

	let { show }: { show: ShowInfo } = $props();
	let metadata = $state<ShowMetadataInfo[]>([]);
	let loading = $state(true);
	let showMetadata = $state(true);

	async function loadMetadata() {
		try {
			if (!auth.token) return;
			metadata = await loadShowMetadata(auth.token, show.id);
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadMetadata();
	});
</script>

<div class="rounded bg-neutral-800 p-4">
	<div class="font-bold">
		{show.name}
	</div>
	{#if loading}
		Loading metadata...
	{:else}
		<div class="flex items-center justify-between">
			<span>
				Metadata: {metadata.length}
			</span>
			<button
				onclick={() => {
					showMetadata = !showMetadata;
				}}
			>
				{#if showMetadata}
					Hide
				{:else}
					Show
				{/if}
			</button>
		</div>
		{#if showMetadata}
			<div class="mt-4 flex flex-col gap-2">
				{#each metadata as m (m.id)}
					<div class="">
						<div class="font-bold">{m.id}</div>
						<div class="font-bold">{m.name}</div>
						<div class="grid grid-cols-2">
							<div>
								<p>{m.description}</p>
								{#if metadata.length > 1}
									<button
										onclick={() => {
											if (!auth.token) return;
											setPrimaryMetadataForShow(auth.token, show.id, m.id)
												.then(() => {
													loadMetadata();
												})
												.catch((e) => {
													const m = (e as Error).message;
													alert(m);
												});
										}}
										class="mt-2 cursor-pointer rounded bg-neutral-700 px-4 py-2 hover:bg-neutral-600"
										>Select as Primary</button
									>
								{/if}
							</div>
							<div>
								<img class="w-full" src={m.medium_image_url} alt={m.name} />
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	{/if}
</div>
