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
		<div>{d.show_id}</div>
	{/each}
</div>
