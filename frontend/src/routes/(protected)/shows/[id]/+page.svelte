<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { page } from '$app/state';
	import { type SeasonInfo, type SeasonListResponse, type ShowInfo } from '$lib/types/export_types';
	import { resolve } from '$app/paths';

	const showId = $derived(page.params.id ?? '');

	let loadingShowData = $state(true);
	let loadingSeasons = $state(true);

	let errorMessage = $state<string | null>(null);

	let showData = $state<ShowInfo | null>(null);
	let seasonData = $state<SeasonInfo[] | null>(null);

	let selectedSeason = $state<SeasonInfo | null>(null);

	async function loadSeasons() {
		try {
			const res = await fetch('/api/seasons/' + showId, {
				headers: {
					Authorization: 'Bearer ' + auth.token
				}
			});
			if (!res.ok) {
				throw new Error(`Failed to load seasons: ${res.status}`);
			}

			const r = (await res.json()) as SeasonListResponse;
			seasonData = r.seasons;

			if (selectedSeason == null && seasonData.length > 0) {
				selectedSeason = seasonData[0];
			}
		} catch (error) {
		} finally {
			loadingSeasons = false;
		}
	}

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
		loadSeasons();
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
		{:else}
			<div class="flex gap-2">
				{#each seasonData as season (season.id)}
					<button
						class={`cursor-pointer ${selectedSeason == season ? 'font-bold' : ''}`}
						onclick={() => (selectedSeason = season)}
					>
						Season {season.number}
					</button>
				{/each}
			</div>
		{/if}
	</div>
</main>
