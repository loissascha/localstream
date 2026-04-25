<script lang="ts">
	import { loadShows } from '$lib/api/shows';
	import { auth } from '$lib/auth.svelte';
	import AdminShowMetadataBlock from '$lib/components/admin/AdminShowMetadataBlock.svelte';
	import { type ShowInfo } from '$lib/types/export_types';

	var shows = $state<ShowInfo[]>([]);
	var loadingShows = $state(true);
	var errorMessage = $state('');
	var hideSingle = $state(true);

	async function loadShowsList() {
		try {
			if (!auth.token) return;
			const data = await loadShows(auth.token);
			shows = data.shows;
		} catch (e) {
			errorMessage = (e as Error).message;
		} finally {
			loadingShows = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadShowsList();
	});
</script>

<div class="mb-4">
	<input id="hidesingle" bind:checked={hideSingle} class="cursor-pointer" type="checkbox" />
	<label for="hidesingle" class="cursor-pointer"> Hide items with single metadata</label>
</div>

{#if errorMessage != ''}
	<p class="text-red-500">{errorMessage}</p>
{/if}

<section class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
	{#each shows as show (show.id)}
		{#if !hideSingle || show.fetch_source == 'multiple' || show.fetch_source == 'empty'}
			<AdminShowMetadataBlock {show} />
		{/if}
	{/each}
</section>
