<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { page } from '$app/state';
	import {
		type EpisodeInfo,
		type EpisodeListResponse,
		type SeasonInfo
	} from '$lib/types/export_types';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { deleteWatchstate, setWatchstateFinished } from '$lib/api/watchstate';
	import { loadSeasonsForShow } from '$lib/api/seasons';
	import CheckIcon from '$lib/icons/CheckIcon.svelte';
	import DOMPurify from 'dompurify';
	import PlusIcon from '$lib/icons/PlusIcon.svelte';
	import SelectCollectionOverlay from '$lib/components/overlays/SelectCollectionOverlay.svelte';
	import { addShowToCollection } from '$lib/api/collections';
	import { shows } from '$lib/shows.svelte';
	import PercentageBar from '$lib/components/ui/PercentageBar.svelte';

	const showId = $derived(page.params.showID ?? '');

	let loadingSeasons = $state(true);
	let loadingEpisodes = $state(true);

	let showAddToCollection = $state(false);

	let errorMessage = $state<string | null>(null);

	let show = $derived.by(() => {
		return shows.shows.find((show) => show.id === showId);
	});
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

	$effect(() => {
		if (!showId || showId == '') return;
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
			{#if show != null && show.medium_image_url != null}
				<div>
					<img alt={show.name} class="max-h-102" src={show.medium_image_url} />
				</div>
			{/if}
		</div>
		<div>
			<h1 class="mb-2 text-3xl">
				{show?.name}
				{#if show && show.year > 0}
					({show.year})
				{/if}
			</h1>
			<div>
				{@html DOMPurify.sanitize(show?.description ?? '')}
			</div>
			<div class="p-4">
				<button
					onclick={() => {
						showAddToCollection = true;
					}}
					class="mt-4 flex cursor-pointer gap-1 rounded bg-neutral-800 px-4 py-2 font-semibold hover:bg-neutral-700"
				>
					<PlusIcon />
					Add to Collection
				</button>
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
				class="relative flex aspect-video shrink-0 flex-col justify-between rounded md:w-54"
			>
				<div>
					{#if episode.medium_image_url != ''}
						<img
							class="aspect-video w-full rounded"
							alt={`Episode ${episode.number}`}
							src={episode.medium_image_url}
						/>
					{:else}
						<div class="aspect-video w-full rounded bg-neutral-800"></div>
					{/if}
				</div>
				<div class="mt-1 w-full truncate font-bold text-neutral-400">
					E{episode.number}: {episode.name}
				</div>
				<div>
					{#if episode.watchstate.percentage > 0 && !episode.watchstate.finished}
						<PercentageBar percentage={episode.watchstate.percentage} />
					{:else}
						<div class="h-2 w-full"></div>
					{/if}
				</div>

				<div class="absolute top-2 right-2 z-10">
					{#if episode.watchstate.finished}
						<button
							class="cursor-pointer rounded-full bg-neutral-800 text-brand"
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
							class="cursor-pointer rounded-full bg-neutral-800 text-neutral-100"
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
			</a>
		{/each}
	</div>
</main>

{#if showAddToCollection}
	<SelectCollectionOverlay
		selectedCollection={(collectionId) => {
			if (auth.token && show) {
				addShowToCollection(auth.token, collectionId, show.id)
					.then(() => {
						showAddToCollection = false;
					})
					.catch((e) => {
						alert((e as Error).message);
					});
			}
		}}
		close={() => {
			showAddToCollection = false;
		}}
	/>
{/if}
