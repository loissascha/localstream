<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { page } from '$app/state';
	import {
		type ShowMetadataInfo,
		type EpisodeInfo,
		type EpisodeListResponse,
		type SeasonInfo,
		type SeasonListResponse,
		type ShowInfo
	} from '$lib/types/export_types';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { setWatchstateFinished } from '$lib/api/watchstate';
	import { loadShowMetadata } from '$lib/api/show_metadata';
	import { loadSeasonsForShow } from '$lib/api/seasons';
	import CheckIcon from '$lib/icons/CheckIcon.svelte';

	const showId = $derived(page.params.showID ?? '');

	let loadingShowData = $state(true);
	let loadingSeasons = $state(true);
	let loadingEpisodes = $state(true);

	let errorMessage = $state<string | null>(null);

	let showData = $state<ShowInfo | null>(null);
	let seasonData = $state<SeasonInfo[] | null>(null);
	let episodeData = $state<EpisodeInfo[] | null>(null);

	let metadata = $state<ShowMetadataInfo | null>(null);

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
			console.log('loaded episodes:', episodeData);
		} catch (error) {
			errorMessage =
				error instanceof Error ? error.message : 'Unknown error while loading show data';
		} finally {
			loadingEpisodes = false;
		}
	}

	async function loadSeasons() {
		try {
			if (!auth.token) return;
			const r = await loadSeasonsForShow(auth.token, showId);
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

	async function loadMetadata() {
		try {
			if (!auth.token) return;
			const data = await loadShowMetadata(auth.token, showId);
			if (data.length == 1) {
				metadata = data[0];
			}
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		}
	}

	$effect(() => {
		if (!showId || showId == '') return;
		loadShows();
		loadSeasons();
		loadMetadata();
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

<main>
	{#if errorMessage}
		<p>{errorMessage}</p>
	{/if}

	<div class="flex flex-col gap-2 md:flex-row">
		<div class="shrink-0">
			{#if metadata != null}
				<div>
					<img alt={metadata.name} class="max-h-102" src={metadata.original_image_url} />
				</div>
			{/if}
		</div>
		<div>
			{#if loadingShowData}
				<p>Loading stuff...</p>
			{:else}
				<h1 class="mb-2 text-3xl">
					{showData?.name}
					{#if showData && showData.year > 0}
						({showData.year})
					{/if}
				</h1>
			{/if}
			{#if metadata != null}
				<div>
					{metadata.description.replaceAll('<p>', '').replaceAll('</p>', '')}
				</div>
			{/if}
		</div>
	</div>

	<div class="my-3">
		{#if loadingSeasons}
			<p>Loading seasons...</p>
		{:else}
			<div class="flex flex-wrap gap-2">
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
	<div class="my-3 grid grid-cols-3 gap-4 py-5 md:flex md:flex-wrap">
		{#each episodeData as episode (episode.id)}
			<a
				href={resolve(
					'/(protected)/(watch)/shows/[showID]/seasons/[seasonID]/episodes/[episodeID]',
					{
						showID: showId,
						seasonID: selectedSeason!.id,
						episodeID: episode.id
					}
				)}
				class="flex aspect-square shrink-0 flex-col justify-between rounded bg-neutral-800 md:w-34"
			>
				<div>
					<div class="flex justify-end px-2 py-1">
						{#if episode.watchstate.finished}
							<div class="cursor-pointer text-brand"><CheckIcon /></div>
						{:else}
							<button
								class="cursor-pointer text-neutral-500"
								onclick={(e) => {
									e.preventDefault();
									e.stopPropagation();
									if (auth.token) {
										setWatchstateFinished(auth.token, episode.id)
											.then((res) => {
												if (selectedSeason) {
													loadEpisodes(selectedSeason.id);
												}
											})
											.catch((e) => {
												const m = (e as Error).message;
												alert(m);
											});
									}
								}}
							>
								<CheckIcon />
							</button>
						{/if}
					</div>
				</div>
				<div class="grow content-center text-center text-2xl font-bold">
					{episode.number}
				</div>
				<div>
					{#if episode.watchstate.percentage > 0 && !episode.watchstate.finished}
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
