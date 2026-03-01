<script lang="ts">
	import { onDestroy } from 'svelte';
	import { page } from '$app/state';

	const videoId = $derived(page.params.id ?? '');
	const streamUrl = $derived(`/api/videos/stream?id=${encodeURIComponent(videoId)}`);

	let videoEl = $state<HTMLVideoElement | null>(null);
	let logTimer: ReturnType<typeof setInterval> | null = null;

	const stopPlaybackLogging = () => {
		if (logTimer !== null) {
			clearInterval(logTimer);
			logTimer = null;
		}
	};

	const logPlaybackStatus = () => {
		if (!videoEl) {
			return;
		}

		const duration = Number.isFinite(videoEl.duration) ? Number(videoEl.duration.toFixed(2)) : 0;
		const position = Number(videoEl.currentTime.toFixed(2));
		const finished = duration > 0 && position >= Math.max(duration - 1, 0);

		console.log({
			videoId,
			positionSeconds: position,
			durationSeconds: duration,
			finished,
			updatedAt: new Date().toISOString()
		});
	};

	const startPlaybackLogging = () => {
		stopPlaybackLogging();
		logPlaybackStatus();
		logTimer = setInterval(logPlaybackStatus, 5000);
	};

	onDestroy(() => {
		stopPlaybackLogging();
	});
</script>

<main class="grid min-h-dvh grid-rows-[auto_1fr] bg-slate-950 text-slate-50">
	<header class="border-b border-slate-400/20 bg-gradient-to-b from-slate-900/95 to-slate-900/70 px-4 py-3.5">
		<a
			class="inline-block rounded-md border border-slate-400/30 px-2.5 py-1.5 text-sm text-slate-300 no-underline hover:border-slate-300/70 hover:text-slate-50"
			href="/"
		>
			Back to library
		</a>
	</header>

	<section class="flex items-center justify-center p-[clamp(0.5rem,1.2vw,1rem)]">
		<!-- svelte-ignore a11y_media_has_caption -->
		<video
			bind:this={videoEl}
			class="h-auto max-h-[calc(100dvh-4.2rem)] w-full rounded-xl bg-black md:h-[min(88dvh,calc(100dvh-4.2rem))] md:w-[min(100%,120rem)]"
			controls
			autoplay
			preload="metadata"
			src={streamUrl}
			onplay={startPlaybackLogging}
			onpause={stopPlaybackLogging}
			onended={stopPlaybackLogging}
		></video>
	</section>
</main>
