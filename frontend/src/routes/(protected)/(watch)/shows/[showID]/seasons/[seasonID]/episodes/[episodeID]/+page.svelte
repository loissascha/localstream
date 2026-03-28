<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { getWatchstateForEpisode, updateWatchstate } from '$lib/api/watchstate';
	import { auth } from '$lib/auth.svelte';
	import { API_URL } from '$lib/consts';
	import { onDestroy } from 'svelte';

	let videoEl = $state<HTMLVideoElement | null>(null);

	const showId = $derived(page.params.showID ?? '');
	const seasonId = $derived(page.params.seasonID ?? '');
	const episodeId = $derived(page.params.episodeID ?? '');
	var loadingWatchstate = $state(true);

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
</script>

<main class="grid min-h-dvh grid-rows-[auto_1fr]">
	<header class="border-b border-neutral-500 bg-neutral-800 px-4 py-3.5">
		<a
			class="inline-block rounded-md border border-slate-400/30 px-2.5 py-1.5 text-sm text-slate-300 no-underline hover:border-slate-300/70 hover:text-slate-50"
			href={resolve('/(protected)')}>Dashboard</a
		>
		<a
			class="inline-block rounded-md border border-slate-400/30 px-2.5 py-1.5 text-sm text-slate-300 no-underline hover:border-slate-300/70 hover:text-slate-50"
			href={resolve('/(protected)/(watch)/shows/[showID]', { showID: showId })}
		>
			Back to Show
		</a>
	</header>

	<section class="flex items-center justify-center p-[clamp(0.5rem,1.2vw,1rem)]">
		<!-- svelte-ignore a11y_media_has_caption -->
		<video
			bind:this={videoEl}
			class="h-auto max-h-[calc(100dvh-4.2rem)] w-full rounded-xl bg-black md:h-[min(88dvh,calc(100dvh-4.2rem))] md:w-[min(100%,120rem)]"
			controls
			preload="metadata"
			src={streamUrl}
			onplay={startPlaybackLogging}
			onpause={stopPlaybackLogging}
			onended={stopPlaybackLogging}
		></video>
	</section>
</main>
