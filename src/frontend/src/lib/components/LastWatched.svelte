<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { listLatestWatchstateByShow } from '$lib/api/watchstate';
	import { auth } from '$lib/auth.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import type { ShowInfo, WatchstateResponse } from '$lib/types/export_types';
	import ItemGrid from './ItemGrid.svelte';
	import ListItemA from './ListItemA.svelte';
	import ShowMetadataSearchOverlay from './overlays/ShowMetadataSearchOverlay.svelte';
	import ShowInfoDisplay from './ShowInfoDisplay.svelte';
	import ContextMenu from './ui/ContextMenu.svelte';

	let data = $state<WatchstateResponse[]>([]);

	let showMetadataOverlayOpen = $state(false);
	let showMetadataShow = $state<ShowInfo | null>(null);

	async function updateData() {
		if (!auth.token) {
			return;
		}
		try {
			const result = await listLatestWatchstateByShow(auth.token);
			data = result;
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.loggedIn) {
			goto(resolve('/(auth)/login'));
			return;
		}
		updateData();
	});
</script>

<section class="my-4">
	{#if data.length > 0}
		<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
			<ChevronRightIcon />
			Continue Shows
		</h2>
		<ItemGrid>
			{#each data as d (d.id)}
				{#if !d.finished}
					<ContextMenu closeOnItemClick={true}>
						<ListItemA
							href={resolve(
								'/(protected)/watch/shows/[showID]/seasons/[seasonID]/episodes/[episodeID]',
								{
									showID: d.show_id,
									seasonID: d.season_id,
									episodeID: d.episode_id
								}
							)}
						>
							<ShowInfoDisplay
								show={d.show_info}
								nameLink
								percentage={d.percentage}
								showPercentage
							/>
							<div>
								<div>S{d.season_info.number}:E{d.episode_info.number}</div>
							</div>
						</ListItemA>
						{#snippet items(closeMenu)}
							<button
								role="menuitem"
								onclick={(e) => {
									e.preventDefault();
									e.stopPropagation();
									closeMenu();
									showMetadataOverlayOpen = true;
									showMetadataShow = d.show_info;
								}}
								class="cursor-pointer px-4 py-2 hover:bg-neutral-700">Update Metadata</button
							>
						{/snippet}
					</ContextMenu>
				{/if}
			{/each}
		</ItemGrid>
	{/if}
</section>

{#if showMetadataOverlayOpen && showMetadataShow}
	<ShowMetadataSearchOverlay
		show={showMetadataShow}
		close={() => {
			showMetadataOverlayOpen = false;
		}}
	/>
{/if}
