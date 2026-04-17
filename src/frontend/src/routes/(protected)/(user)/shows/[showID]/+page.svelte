<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { page } from '$app/state';
	import {
		type EpisodeInfo,
		type EpisodeListResponse,
		type SeasonInfo,
		type ShowInfo
	} from '$lib/types/export_types';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { deleteWatchstate, setWatchstateFinished } from '$lib/api/watchstate';
	import { loadSeasonsForShow } from '$lib/api/seasons';
	import CheckIcon from '$lib/icons/CheckIcon.svelte';
	import DOMPurify from 'dompurify';

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

	async function loadShowData() {
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
		loadShowData();
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

<main>
	{#if errorMessage}
		<p>{errorMessage}</p>
	{/if}

	<div class="flex flex-col gap-2 md:flex-row">
		<div class="shrink-0">
			{#if showData != null && showData.medium_image_url != null}
				<div>
					<img alt={showData.name} class="max-h-102" src={showData.medium_image_url} />
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
			<div>
				{@html DOMPurify.sanitize(showData?.description ?? '')}
			</div>
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
				href={resolve('/(protected)/watch/shows/[showID]/seasons/[seasonID]/episodes/[episodeID]', {
					showID: showId,
					seasonID: selectedSeason!.id,
					episodeID: episode.id
				})}
				class="flex aspect-square shrink-0 flex-col justify-between rounded bg-neutral-800 md:w-34"
			>
				<div>
					<div class="flex justify-end px-2 py-1">
						{#if episode.watchstate.finished}
							<button
								class="cursor-pointer text-brand"
								onclick={(e) => {
									e.preventDefault();
									e.stopPropagation();
									if (auth.token && selectedSeason) {
										deleteWatchstate(auth.token, episode.id)
											.then(() => {
												if (selectedSeason) {
													loadEpisodes(selectedSeason.id);
												}
											})
											.catch((e) => {
												const m = (e as Error).message;
												alert(m);
											});
									}
								}}><CheckIcon /></button
							>
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
							class={`h-2 bg-brand text-sm`}
						></div>
					{:else}
						<div class="h-2 w-full"></div>
					{/if}
				</div>
			</a>
		{/each}
	</div>
</main>
