<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { getNextEpisode } from '$lib/api/episode';
	import { getWatchstateForEpisode, updateWatchstate } from '$lib/api/watchstate';
	import { auth } from '$lib/auth.svelte';
	import { API_URL } from '$lib/consts';
	import ChevronLeftIcon from '$lib/icons/ChevronLeftIcon.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	import { type EpisodeInfo } from '$lib/types/export_types';
	import { onDestroy } from 'svelte';

	let videoEl = $state<HTMLVideoElement | null>(null);

	const showId = $derived(page.params.showID ?? '');
	const seasonId = $derived(page.params.seasonID ?? '');
	const episodeId = $derived(page.params.episodeID ?? '');
	var loadingWatchstate = $state(true);
	var almostDone = $state(false);
	var nextEpisode = $state<EpisodeInfo | null>(null);

	const streamUrl = $derived(
		API_URL +
			`/api/episodes/stream?id=${encodeURIComponent(episodeId)}&token=${encodeURIComponent(auth.token ? auth.token : '')}`
	);
	let logTimer: ReturnType<typeof setInterval> | null = null;

	function stopPlaybackLogging() {
		if (logTimer !== null) {
			clearInterval(logTimer);
			logTimer = null;
		}
	}

	async function logPlaybackStatus() {
		if (!videoEl) {
			return;
		}
		if (loadingWatchstate) {
			return;
		}

		const duration = Number.isFinite(videoEl.duration) ? Number(videoEl.duration.toFixed(2)) : 0;
		const position = Number(videoEl.currentTime.toFixed(2));
		const finished = duration > 0 && position >= Math.max(duration - 10, 0);
		almostDone = duration > 0 && position >= Math.max(duration - 120, 0);

		if (auth.token) {
			try {
				await updateWatchstate(auth.token, {
					episode_id: episodeId,
					season_id: seasonId,
					show_id: showId,
					position: position,
					duration: duration,
					finished: finished
				});
			} catch (e) {
				console.error(e);
			}
		}

		// 	console.log({
		// 		userToken: auth.token,
		// 		showId: showId,
		// 		seasonId: seasonId,
		// 		episodeId: episodeId,
		// 		positionSeconds: position,
		// 		durationSeconds: duration,
		// 		finished,
		// 		updatedAt: new Date().toISOString()
		// 	});
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
		if (videoEl == null) return;
		loadingWatchstate = true;
		getWatchstateForEpisode(auth.token, episodeId)
			.then((res) => {
				console.log('watchstate res:', res);
				if (!res.finished) {
					videoEl!.currentTime = res.position;
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
		if (!auth.token) {
			return;
		}
		if (!loadingWatchstate) return;
		if (!almostDone) return;
		if (nextEpisode != null) return;
		getNextEpisode(auth.token, episodeId)
			.then((r) => {
				nextEpisode = r;
			})
			.catch((e) => {
				const m = (e as Error).message;
				alert(m);
			});
	});
</script>

<main class="grid min-h-dvh grid-rows-[auto_1fr]">
	<header class="flex items-center gap-2 bg-neutral-900 px-4 py-3.5">
		<a class="p-2 text-slate-300 no-underline hover:text-white" href={resolve('/(protected)')}>
			<HomeIcon />
		</a>
		<a
			class="p-2 text-slate-300 no-underline hover:text-white"
			href={resolve('/(protected)/(watch)/shows/[showID]', { showID: showId })}
		>
			<ChevronLeftIcon />
		</a>
	</header>

	<section class="flex items-center justify-center bg-black">
		<!-- svelte-ignore a11y_media_has_caption -->
		<video
			bind:this={videoEl}
			class="h-full w-full bg-black"
			controls
			preload="metadata"
			src={streamUrl}
			onplay={startPlaybackLogging}
			onpause={stopPlaybackLogging}
			onended={stopPlaybackLogging}
		></video>
	</section>

	{#if almostDone && nextEpisode != null}
		<a
			href={resolve('/(protected)/(watch)/shows/[showID]/seasons/[seasonID]/episodes/[episodeID]', {
				showID: showId,
				seasonID: nextEpisode.season_id,
				episodeID: nextEpisode.id
			})}
			class="fixed right-8 bottom-12 flex h-10 w-60 items-center justify-center rounded bg-neutral-500 border border-neutral-600"
		>
			Next Episode <ChevronRightIcon />
		</a>
	{/if}
</main>
