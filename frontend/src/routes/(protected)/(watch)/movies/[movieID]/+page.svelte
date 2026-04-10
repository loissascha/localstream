<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { getWatchstateForMovie, updateWatchstateMovie } from '$lib/api/watchstate_movie';
	import { auth } from '$lib/auth.svelte';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	import { onDestroy } from 'svelte';

	let videoEl = $state<HTMLVideoElement | null>(null);

	const movieId = $derived(page.params.movieID ?? '');
	var loadingWatchstate = $state(true);

	const streamUrl = $derived(
		`/api/movies/stream?id=${encodeURIComponent(movieId)}&token=${encodeURIComponent(auth.token ? auth.token : '')}`
	);
	let logTimer: ReturnType<typeof setInterval> | null = null;

	function stopPlaybackLogging() {
		if (logTimer !== null) {
			clearInterval(logTimer);
			logTimer = null;
			console.log('playback logging stopped');
		}
	}

	async function logPlaybackStatus() {
		console.log('log playback status 1');
		if (!videoEl) {
			return;
		}
		console.log('log playback status 2');
		if (loadingWatchstate) {
			return;
		}

		const duration = Number.isFinite(videoEl.duration) ? Number(videoEl.duration.toFixed(2)) : 0;
		const position = Number(videoEl.currentTime.toFixed(2));
		const finished = duration > 0 && position >= Math.max(duration - 10, 0);

		if (auth.token) {
			try {
				console.log('log watchsatte');
				await updateWatchstateMovie(auth.token, {
					movie_id: movieId,
					position: position,
					duration: duration,
					finished: finished
				});
			} catch (e) {
				console.error(e);
			}
		}
	}

	function startPlaybackLogging() {
		console.log('start playback logging');
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
		getWatchstateForMovie(auth.token, movieId)
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
</script>

<main class="grid h-dvh grid-rows-[auto_1fr] overflow-hidden">
	<header class="flex items-center gap-2 bg-neutral-900 px-4 py-3.5">
		<a class="p-2 text-slate-300 no-underline hover:text-white" href={resolve('/(protected)')}>
			<HomeIcon />
		</a>
	</header>

	<section class="min-h-0 bg-orange-100">
		<!-- svelte-ignore a11y_media_has_caption -->
		<div class="h-full w-full bg-green-100">
			<video
				bind:this={videoEl}
				class="h-full w-full bg-black object-contain"
				controls
				preload="metadata"
				src={streamUrl}
				onplay={startPlaybackLogging}
				onpause={stopPlaybackLogging}
				onended={stopPlaybackLogging}
			></video>
		</div>
	</section>
</main>
