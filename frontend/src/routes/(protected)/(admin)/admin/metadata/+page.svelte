<script lang="ts">
	import { loadShows } from '$lib/api/shows';
	import { auth } from '$lib/auth.svelte';
	import AdminShowMetadataBlock from '$lib/components/admin/AdminShowMetadataBlock.svelte';
	import type { ShowInfo } from '$lib/types/export_types';

	var shows = $state<ShowInfo[]>([]);
	var loadingShows = $state(true);
	var errorMessage = $state('');

	async function loadShowsList() {
		try {
			if (!auth.token) return;
			const data = await loadShows(auth.token);
			shows = data.shows;
		} catch (e) {
			errorMessage = (e as Error).message;
		} finally {
			loadingShows = true;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadShowsList();
	});
</script>

{#if errorMessage != ''}
	<p class="text-red-500">{errorMessage}</p>
{/if}

<div class="grid grid-cols-3 gap-4">
	{#each shows as show}
		<AdminShowMetadataBlock {show} />
	{/each}
</div>
