<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { listLatestWatchstateByShow } from '$lib/api/watchstate';
	import { auth } from '$lib/auth.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import { type WatchstateResponse } from '$lib/types/export_types';

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

<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
	<ChevronRightIcon />
	Continue
</h2>
<div class="flex gap-4">
	{#each data as d}
		<a
			href={resolve('/(protected)/(watch)/shows/[showID]/seasons/[seasonID]/episodes/[episodeID]', {
				showID: d.show_id,
				seasonID: d.season_id,
				episodeID: d.episode_id
			})}
			class="flex w-60 cursor-pointer flex-col justify-between gap-2 rounded bg-neutral-800 p-4 hover:bg-neutral-700"
		>
			<div class="font-bold">{d.show_info.name}</div>
			<div>
				<div>Season: {d.season_info.number}</div>
				<div>Episode: {d.episode_info.number}</div>
				<div class="bg-neutral-600">
					<div style={`width: ${d.percentage}%;`} class={`h-2 bg-blue-300 text-sm`}></div>
				</div>
			</div>
		</a>
	{/each}
</div>
