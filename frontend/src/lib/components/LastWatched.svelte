<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { listLatestWatchstateByShow } from '$lib/api/watchstate';
	import { auth } from '$lib/auth.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import { type WatchstateResponse } from '$lib/types/export_types';
	import ItemGrid from './ItemGrid.svelte';
	import ListItemA from './ListItemA.svelte';
	import ShowInfoDisplay from './ShowInfoDisplay.svelte';

	let data = $state<WatchstateResponse[]>([]);

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
	<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
		<ChevronRightIcon />
		Continue Shows
	</h2>
	<ItemGrid>
		{#each data as d (d.id)}
			<ListItemA
				href={resolve(
					'/(protected)/(watch)/shows/[showID]/seasons/[seasonID]/episodes/[episodeID]',
					{
						showID: d.show_id,
						seasonID: d.season_id,
						episodeID: d.episode_id
					}
				)}
			>
				<ShowInfoDisplay show={d.show_info} />
				<div>
					<div>S{d.season_info.number} E{d.episode_info.number}</div>
					<div class="bg-neutral-600">
						<div style={`width: ${d.percentage}%;`} class={`h-2 bg-brand text-sm`}></div>
					</div>
				</div>
			</ListItemA>
		{/each}
	</ItemGrid>
</section>
