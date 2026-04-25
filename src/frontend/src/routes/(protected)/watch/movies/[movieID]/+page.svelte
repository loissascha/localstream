<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { getWatchstateForMovie } from '$lib/api/watchstate_movie';
	import { auth } from '$lib/auth.svelte';
	import VideoPlayer from '$lib/components/VideoPlayer.svelte';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	import { setMovieWatchstate } from '$lib/movies.svelte';
	import { onDestroy } from 'svelte';

	const movieId = $derived(page.params.movieID ?? '');
	var loadingWatchstate = $state(true);

	let duration = $state(0);
	let currentTime = $state(0);

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
		console.log('log playback status 2');
		if (loadingWatchstate) {
			return;
		}

		const position = Number(currentTime.toFixed(2));

		try {
			console.log('log watchsatte');
			await setMovieWatchstate(movieId, position, duration);
		} catch (e) {
			console.error(e);
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
		loadingWatchstate = true;
		getWatchstateForMovie(auth.token, movieId)
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
			{/snippet}
		</VideoPlayer>
	</section>
</main>
