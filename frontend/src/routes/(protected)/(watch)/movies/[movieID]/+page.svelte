<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { auth } from '$lib/auth.svelte';
	import { API_URL } from '$lib/consts';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	// import { onDestroy } from 'svelte';

	let videoEl = $state<HTMLVideoElement | null>(null);

	const movieId = $derived(page.params.movieID ?? '');
	// var loadingWatchstate = $state(true);

	const streamUrl = $derived(
		API_URL +
			`/api/movies/stream?id=${encodeURIComponent(movieId)}&token=${encodeURIComponent(auth.token ? auth.token : '')}`
	);
	// let logTimer: ReturnType<typeof setInterval> | null = null;

	 // function stopPlaybackLogging() {
	 // 	if (logTimer !== null) {
	 // 		clearInterval(logTimer);
	 // 		logTimer = null;
	 // 	}
	 // }
	
	 // async function logPlaybackStatus() {
	 // 	if (!videoEl) {
	 // 		return;
	 // 	}
	 // 	if (loadingWatchstate) {
	 // 		return;
	 // 	}
	 //
	 // 	const duration = Number.isFinite(videoEl.duration) ? Number(videoEl.duration.toFixed(2)) : 0;
	 // 	const position = Number(videoEl.currentTime.toFixed(2));
	 // 	const finished = duration > 0 && position >= Math.max(duration - 10, 0);
	 // 	// almostDone = duration > 0 && position >= Math.max(duration - 120, 0);
	 //
	 // 	if (auth.token) {
	 // 		try {
	 // 		console.log("log watchsatte")
	 // 		} catch (e) {
	 // 			console.error(e);
	 // 		}
	 // 	}
	 // }
	 //
	 // function startPlaybackLogging() {
	 // 	stopPlaybackLogging();
	 // 	logPlaybackStatus();
	 // 	logTimer = setInterval(logPlaybackStatus, 5000);
	 //
	 //
	 // onDestroy(() => {
	 // 	stopPlaybackLogging();
	 // });
	 //
	 // $effect(() => {
	 // 	if (!auth.token) {
	 // 		return;
	 // 	}
	 // 	if (videoEl == null) return;
	 // 	loadingWatchstate = false;
	 // });
</script>

<main class="grid min-h-dvh grid-rows-[auto_1fr]">
	<header class="flex items-center gap-2 bg-neutral-900 px-4 py-3.5">
		<a class="p-2 text-slate-300 no-underline hover:text-white" href={resolve('/(protected)')}>
			<HomeIcon />
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
</main>
