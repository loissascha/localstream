<script lang="ts">
	import { loadShows } from '$lib/api/shows';
	import { auth } from '$lib/auth.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import ShowListItem from '$lib/components/ShowListItem.svelte';
	import LibraryIcon from '$lib/icons/LibraryIcon.svelte';
	import type { ShowInfo } from '$lib/types/export_types';

	let shows = $state<ShowInfo[]>([]);
	let loading = $state(true);
	let errorMessage = $state('');

	async function loadShowsList() {
		try {
			if (!auth.token) return;
			const data = await loadShows(auth.token);
			shows = data.shows;
			shows.sort((a, b) => {
				if (a.name > b.name) return 1;
				if (a.name < b.name) return -1;
				return 0;
			});
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		loadShowsList();
	});
</script>

<main>
	{#if errorMessage}
		<p class="text-red-700">{errorMessage}</p>
	{/if}

	{#if loading}
		<p>Loading shows...</p>
	{:else}
		<section class="my-8">
			<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
				<LibraryIcon />
				Shows
			</h2>
			<ItemGrid>
				{#each shows as show (show.id)}
					<ShowListItem {show} />
				{/each}
			</ItemGrid>
		</section>
	{/if}
</main>
