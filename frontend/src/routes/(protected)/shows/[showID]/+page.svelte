<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { page } from '$app/state';
	import {
		type EpisodeInfo,
		type EpisodeListResponse,
		type SeasonInfo,
		type SeasonListResponse,
		type ShowInfo
	} from '$lib/types/export_types';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';

	const showId = $derived(page.params.showID ?? '');

	let loadingShowData = $state(true);
	let loadingSeasons = $state(true);
	let loadingEpisodes = $state(true);

	let errorMessage = $state<string | null>(null);

	let showData = $state<ShowInfo | null>(null);
	let seasonData = $state<SeasonInfo[] | null>(null);
	let episodeData = $state<EpisodeInfo[] | null>(null);

	let selectedSeason = $state<SeasonInfo | null>(null);

	async function loadEpisodes(seasonID: string) {
		try {
			const res = await fetch('/api/episodes/' + seasonID, {
				headers: {
					Authorization: 'Bearer ' + auth.token
				}
			});
			if (!res.ok) {
				throw new Error(`Failed to load episodes: ${res.status}`);
			}

			const r = (await res.json()) as EpisodeListResponse;
			episodeData = r.episodes;
		} catch (error) {
			errorMessage =
				error instanceof Error ? error.message : 'Unknown error while loading show data';
		} finally {
			loadingEpisodes = false;
		}
	}

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
			errorMessage =
				error instanceof Error ? error.message : 'Unknown error while loading show data';
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
		if (!auth.initialized) return;
		if (!auth.loggedIn) {
			goto(resolve('/(auth)/login'));
			return;
		}
		if (!showId || showId == '') return;
		loadShows();
		loadSeasons();
	});

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.loggedIn) {
			goto(resolve('/(auth)/login'));
			return;
		}
		if (!selectedSeason) return;
		loadEpisodes(selectedSeason.id);
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
		<h1 class="text-3xl">
			{showData?.name}
			{#if showData && showData.year > 0}
				({showData.year})
			{/if}
		</h1>
		<p>{showData?.description}</p>
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
	<div class="my-3 flex gap-4 overflow-y-scroll">
		{#each episodeData as episode (episode.id)}
			<a
				href={resolve('/(protected)/shows/[showID]/seasons/[seasonID]/episodes/[episodeID]', {
					showID: showId,
					seasonID: selectedSeason!.id,
					episodeID: episode.id
				})}
				class="flex h-34 w-34 shrink-0 flex-col justify-between rounded bg-neutral-800"
			>
				<div class="grow content-center text-center text-2xl font-bold">
					{episode.number}
				</div>
				<div>
					{#if episode.watchstate.position > 0}
						<div
							style={`width: ${episode.watchstate.percentage}%;`}
							class={`h-2 bg-blue-300 text-sm`}
						></div>
					{:else}
						<div class="h-2 w-full"></div>
					{/if}
				</div>
			</a>
		{/each}
	</div>
</main>
