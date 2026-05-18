<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { loadMovieSubtitles } from '$lib/api/movie_subtitles';
	import { getMovie } from '$lib/api/movies';
	import { getWatchstateForMovie } from '$lib/api/watchstate_movie';
	import { auth } from '$lib/auth.svelte';
	import MovieSubtitleSearchOverlay from '$lib/components/overlays/MovieSubtitleSearchOverlay.svelte';
	import VideoPlayer from '$lib/components/VideoPlayer.svelte';
	import ChevronLeftIcon from '$lib/icons/ChevronLeftIcon.svelte';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	import { setMovieWatchstate } from '$lib/movies.svelte';
	import type { SubtitleInfo, MovieInfo } from '$lib/types/export_types';
	import { onDestroy } from 'svelte';

	const movieId = $derived(page.params.movieID ?? '');

	let movie = $state<MovieInfo | null>(null);
	let loadingWatchstate = $state(true);
	let subtitles = $state<SubtitleInfo[]>([]);
	let subtitleoverlayopen = $state(false);

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
		getMovie(auth.token, movieId)
			.then((res) => {
				movie = res;
			})
			.catch((e) => {
				const m = (e as Error).message;
				alert(m);
			});
		reloadSubtitles();
	});

	function reloadSubtitles() {
		if (!auth.token) return;
		loadMovieSubtitles(auth.token, movieId)
			.then((r) => {
				subtitles = r;
			})
			.catch((e) => {
				const m = (e as Error).message;
				alert(m);
			});
	}
</script>

<main class="grid h-dvh grid-rows-[1fr] overflow-hidden">
	<section class="min-h-0">
		<VideoPlayer
			href={streamUrl}
			onplay={startPlaybackLogging}
			onpause={stopPlaybackLogging}
			onended={stopPlaybackLogging}
			{subtitles}
			bind:currentTime
			bind:duration
		>
			{#snippet topbar()}
				<a class="p-2 text-slate-300 no-underline hover:text-white" href={resolve('/(protected)')}>
					<HomeIcon />
				</a>
				<a
					class="p-2 text-slate-300 no-underline hover:text-white"
					href={resolve('/(protected)/(user)/movies/[movieID]', {
						movieID: movieId
					})}
				>
					<ChevronLeftIcon />
				</a>
				<span>{movie?.name}</span>
			{/snippet}
			{#snippet bottomrightextensions()}
				<button
					class="cursor-pointer"
					onclick={() => {
						subtitleoverlayopen = true;
					}}
				>
					<span class="rounded-full px-2 py-1 text-xs font-medium text-white/85">CC</span>
				</button>
			{/snippet}
		</VideoPlayer>
	</section>
</main>

{#if subtitleoverlayopen && movie != null}
	<MovieSubtitleSearchOverlay
		close={() => {
			subtitleoverlayopen = false;
			reloadSubtitles();
		}}
		{movie}
	/>
{/if}
