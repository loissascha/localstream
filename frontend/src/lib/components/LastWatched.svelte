<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { listLatestWatchstateByShow } from '$lib/api/watchstate';
	import { auth } from '$lib/auth.svelte';
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

<h2>Last Watched</h2>
<div class="flex gap-4">
	{#each data as d}
		<a
			href={resolve('/(protected)/shows/[showID]/seasons/[seasonID]/episodes/[episodeID]', {
				showID: d.show_id,
				seasonID: d.season_id,
				episodeID: d.episode_id
			})}
			class="cursor-pointer rounded border border-neutral-600 bg-neutral-800 p-4"
		>
			<div class="font-bold">{d.show_info.name}</div>
			<div>Season: {d.season_info.number}</div>
			<div>Episode: {d.episode_info.number}</div>
		</a>
	{/each}
</div>
