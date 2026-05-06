<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { getEpisodeDetails, getEpisodeMetadata, getNextEpisode } from '$lib/api/episode';
	import { getSeasonDetails } from '$lib/api/seasons';
	import { getWatchstateForEpisode } from '$lib/api/watchstate';
	import { auth } from '$lib/auth.svelte';
	import VideoPlayer from '$lib/components/VideoPlayer.svelte';
	import ChevronLeftIcon from '$lib/icons/ChevronLeftIcon.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	import { setShowWatchstate } from '$lib/shows.svelte';
	import {
		type SeasonInfo,
		type EpisodeInfo,
		type EpisodeMetadataInfo
	} from '$lib/types/export_types';
	import { onDestroy } from 'svelte';

	const showId = $derived(page.params.showID ?? '');
	const seasonId = $derived(page.params.seasonID ?? '');
	const episodeId = $derived(page.params.episodeID ?? '');

	var loadingWatchstate = $state(true);
	var almostDone = $state(false);
	let duration = $state(0);
	let currentTime = $state(0);

	var episodeMetadataDetails = $state<EpisodeMetadataInfo | null>(null);
	var episodeDetails = $state<EpisodeInfo | null>(null);
	var seasonDetails = $state<SeasonInfo | null>(null);
	var nextEpisode = $state<EpisodeInfo | null>(null);

	const streamUrl = $derived(
		`/api/episodes/stream?id=${encodeURIComponent(episodeId)}&token=${encodeURIComponent(auth.token ? auth.token : '')}`
	);
	let logTimer: ReturnType<typeof setInterval> | null = null;

	async function loadEpisodeDetails() {
		if (!auth.token) return;
		const data = await getEpisodeDetails(auth.token, episodeId);
		episodeDetails = data;
	}

	async function loadSeasonDetails() {
		if (!auth.token) return;
		const data = await getSeasonDetails(auth.token, seasonId);
		seasonDetails = data;
	}

	function stopPlaybackLogging() {
		if (logTimer !== null) {
			clearInterval(logTimer);
			logTimer = null;
		}
	}

	async function logPlaybackStatus() {
		if (loadingWatchstate) {
			return;
		}

		const position = Number(currentTime.toFixed(2));
		const normalizedDuration = Number.isFinite(duration) ? Number(duration.toFixed(2)) : 0;
		const almostDoneStatus =
			normalizedDuration > 0 && position >= Math.max(normalizedDuration - 180, 0);
		console.log('almost done status:', almostDoneStatus);
		almostDone = almostDoneStatus;

		await setShowWatchstate(showId, seasonId, episodeId, position, duration);
	}

	function startPlaybackLogging() {
		stopPlaybackLogging();
		logPlaybackStatus();
		logTimer = setInterval(logPlaybackStatus, 5000);
	}

	onDestroy(() => {
		stopPlaybackLogging();
	});

	$effect(() => {
		if (!auth.token) {
			return;
		}
		loadingWatchstate = true;
		getWatchstateForEpisode(auth.token, episodeId)
			.then((res) => {
				console.log('watchstate res:', res);
				if (!res.finished) {
					currentTime = res.position;
					loadingWatchstate = false;
				} else {
					alert('already watched');
					loadingWatchstate = false;
				}
			})
			.catch((e) => {
				const m = (e as Error).message;
				if (m == '404') {
					// watchstate does not exist -> user didn't watch this yet so it's okay
					loadingWatchstate = false;
				} else {
					alert(m);
				}
			});
	});

	$effect(() => {
		void episodeId;
		if (!auth.initialized) return;
		if (!auth.token) return;
		loadEpisodeDetails();
	});

	$effect(() => {
		void seasonId;
		if (!auth.initialized) return;
		if (!auth.token) return;
		loadSeasonDetails();
	});

	$effect(() => {
		void episodeId;
		almostDone = false;
		nextEpisode = null;
	});

	$effect(() => {
		void episodeDetails;
		void seasonDetails;
		void showId;
		if (!auth.initialized) return;
		if (!auth.token) return;
		if (!seasonDetails) return;
		if (!episodeDetails) return;
		getEpisodeMetadata(auth.token, showId, seasonDetails.number, episodeDetails.number)
			.then((r) => {
				episodeMetadataDetails = r;
			})
			.catch((e) => {
				const m = (e as Error).message;
				alert(m);
			});
	});

	$effect(() => {
		if (!auth.token) {
			return;
		}
		if (loadingWatchstate) return;
		if (!almostDone) return;
		if (nextEpisode != null) return;

		console.log('loading next episode!');
		getNextEpisode(auth.token, episodeId)
			.then((r) => {
				console.log('next episode res:', r);
				nextEpisode = r;
			})
			.catch((e) => {
				const m = (e as Error).message;
				alert(m);
			});
	});
</script>

<main class="grid h-dvh grid-rows-[1fr] overflow-hidden">
	<section class="min-h-0">
		<VideoPlayer
			href={streamUrl}
			onplay={startPlaybackLogging}
			onpause={stopPlaybackLogging}
			onended={stopPlaybackLogging}
			bind:currentTime
			bind:duration
		>
			{#snippet topbar()}
				<a class="p-2 text-slate-300 no-underline hover:text-white" href={resolve('/(protected)')}>
					<HomeIcon />
				</a>
				<a
					class="p-2 text-slate-300 no-underline hover:text-white"
					href={resolve('/(protected)/(user)/shows/[showID]', {
						showID: showId
					}) +
						'?season=' +
						seasonId}
				>
					<ChevronLeftIcon />
				</a>
				<span>
					S{seasonDetails?.number}:E{episodeDetails?.number} - {episodeMetadataDetails?.name}
				</span>
			{/snippet}
			{#snippet overlay()}
				{#if almostDone && nextEpisode != null}
					<a
						href={resolve(
							'/(protected)/watch/shows/[showID]/seasons/[seasonID]/episodes/[episodeID]',
							{
								showID: showId,
								seasonID: nextEpisode.season_id,
								episodeID: nextEpisode.id
							}
						)}
						class="flex h-10 w-60 items-center justify-center gap-2 rounded border border-white/20 bg-black/70 px-4 text-sm text-white backdrop-blur-sm transition hover:bg-black/85"
					>
						Next Episode <ChevronRightIcon />
					</a>
				{/if}
			{/snippet}
		</VideoPlayer>
	</section>
</main>
