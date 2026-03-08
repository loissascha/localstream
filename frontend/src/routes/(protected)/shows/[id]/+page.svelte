<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { page } from '$app/state';
	import { type ShowInfo } from '$lib/types/export_types';
	import { resolve } from '$app/paths';

	const showId = $derived(page.params.id ?? '');

	let loadingShowData = $state(true);
	let loadingSeasons = $state(true);

	let errorMessage = $state<string | null>(null);

	let showData = $state<ShowInfo | null>(null);

	async function loadShows() {
		try {
			const res = await fetch('/api/show/' + showId, {
				headers: {
					Authorization: 'Bearer ' + auth.token
				}
			});
			if (!res.ok) {
				throw new Error(`Failed to load shows: ${res.status}`);
			}

			showData = (await res.json()) as ShowInfo;
		} catch (error) {
			errorMessage =
				error instanceof Error ? error.message : 'Unknown error while loading show data';
		} finally {
			loadingShowData = false;
		}
	}

	$effect(() => {
		if (!showId || showId == '') return;
		loadShows();
	});
</script>

<main class="px-5 py-5">
	<div class="mb-4">
		<a class="cursor-pointer" href={resolve('/(protected)')}>Go Back</a>
	</div>
	{#if errorMessage}
		<p>{errorMessage}</p>
	{/if}

	{#if loadingShowData}
		<p>Loading stuff...</p>
	{:else}
		<h1 class="text-3xl">{showData?.name}</h1>
		<p>Description</p>
	{/if}

	<div class="my-3">
		{#if loadingSeasons}
			<p>Loading seasons...</p>
		{/if}
	</div>
</main>
