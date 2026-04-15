<script lang="ts">
	import { onDestroy } from 'svelte';
	import FullscreenIcon from '$lib/icons/FullscreenIcon.svelte';
	import MuteIcon from '$lib/icons/MuteIcon.svelte';
	import PauseIcon from '$lib/icons/PauseIcon.svelte';
	import PlayIcon from '$lib/icons/PlayIcon.svelte';
	import SkipNextIcon from '$lib/icons/SkipNextIcon.svelte';
	import SkipPreviousIcon from '$lib/icons/SkipPreviousIcon.svelte';
	import VolumeIcon from '$lib/icons/VolumeIcon.svelte';

	interface Props {
		href: string;
		duration?: number;
		currentTime?: number;
		onplay?: () => void;
		onpause?: () => void;
		onended?: () => void;
	}

	let {
		href,
		onplay,
		onpause,
		onended,
		duration = $bindable(0),
		currentTime = $bindable(0)
	}: Props = $props();

	let containerEl = $state<HTMLDivElement | null>(null);
	let videoEl = $state<HTMLVideoElement | null>(null);
	let paused = $state(true);
	let muted = $state(false);
	let volume = $state(1);
	let isFullscreen = $state(false);

	function syncState() {
		if (!videoEl) return;
		currentTime = videoEl.currentTime;
		duration = Number.isFinite(videoEl.duration) ? videoEl.duration : 0;
		paused = videoEl.paused;
		muted = videoEl.muted;
		volume = videoEl.volume;
	}

	function formatTime(value: number) {
		if (!Number.isFinite(value) || value < 0) return '0:00';
		const totalSeconds = Math.floor(value);
		const hours = Math.floor(totalSeconds / 3600);
		const minutes = Math.floor((totalSeconds % 3600) / 60);
		const seconds = totalSeconds % 60;

		if (hours > 0) {
			return `${hours}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
		}

		return `${minutes}:${String(seconds).padStart(2, '0')}`;
	}

	async function togglePlay() {
		if (!videoEl) return;
		if (videoEl.paused) {
			await videoEl.play();
			return;
		}
		videoEl.pause();
	}

	function seekTo(value: number) {
		if (!videoEl) return;
		const boundedValue = Math.min(Math.max(value, 0), duration || 0);
		videoEl.currentTime = boundedValue;
		currentTime = boundedValue;
	}

	function seekBy(seconds: number) {
		seekTo(currentTime + seconds);
	}

	function setVolume(value: number) {
		if (!videoEl) return;
		const boundedValue = Math.min(Math.max(value, 0), 1);
		videoEl.volume = boundedValue;
		videoEl.muted = boundedValue === 0;
		volume = boundedValue;
		muted = videoEl.muted;
	}

	function toggleMute() {
		if (!videoEl) return;
		videoEl.muted = !videoEl.muted;
		muted = videoEl.muted;
	}

	async function toggleFullscreen() {
		if (!containerEl) return;
		if (document.fullscreenElement === containerEl) {
			await document.exitFullscreen();
			return;
		}
		await containerEl.requestFullscreen();
	}

	function handlePlay() {
		syncState();
		onplay?.();
	}

	function handlePause() {
		syncState();
		onpause?.();
	}

	function handleEnded() {
		syncState();
		onended?.();
	}

	function handleFullscreenChange() {
		isFullscreen = document.fullscreenElement === containerEl;
	}

	$effect(() => {
		if (videoEl && Math.abs(videoEl.currentTime - currentTime) > 0.25) {
			videoEl.currentTime = currentTime;
		}
	});

	$effect(() => {
		document.addEventListener('fullscreenchange', handleFullscreenChange);

		return () => {
			document.removeEventListener('fullscreenchange', handleFullscreenChange);
		};
	});

	onDestroy(() => {
		if (document.fullscreenElement === containerEl) {
			document.exitFullscreen().catch(() => {});
		}
	});
</script>

<!-- svelte-ignore a11y_media_has_caption -->
<div bind:this={containerEl} class="relative h-full w-full overflow-hidden bg-black">
	<video
		bind:this={videoEl}
		class="h-full w-full bg-black object-contain"
		preload="metadata"
		playsinline
		src={href}
		onclick={togglePlay}
		onplay={handlePlay}
		onpause={handlePause}
		onended={handleEnded}
		onloadedmetadata={syncState}
		ondurationchange={syncState}
		onratechange={syncState}
		onseeked={syncState}
		onvolumechange={syncState}
		ontimeupdate={syncState}
	></video>

	{#if paused}
		<div class="pointer-events-none absolute inset-0 flex items-center justify-center bg-black/20">
			<button
				type="button"
				class="pointer-events-auto rounded-full bg-white/12 p-5 text-white backdrop-blur-sm transition hover:bg-white/20"
				onclick={togglePlay}
				aria-label="Play video"
			>
				<PlayIcon size={40} />
			</button>
		</div>
	{/if}

	<div
		class="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/90 via-black/55 to-transparent px-4 pt-10 pb-4 text-white"
	>
		<div class="mb-3 flex items-center gap-3 text-xs text-white/80">
			<span class="w-12 text-right tabular-nums">{formatTime(currentTime)}</span>
			<input
				type="range"
				min="0"
				max={duration || 0}
				step="0.1"
				value={currentTime}
				oninput={(event) => seekTo(Number((event.currentTarget as HTMLInputElement).value))}
				class="h-1 w-full cursor-pointer accent-white"
				aria-label="Seek"
			/>
			<span class="w-12 tabular-nums">{formatTime(duration)}</span>
		</div>

		<div class="flex flex-wrap items-center gap-2 sm:flex-nowrap">
			<button
				type="button"
				class="rounded-full p-2 text-white transition hover:bg-white/10"
				onclick={togglePlay}
				aria-label={paused ? 'Play video' : 'Pause video'}
			>
				{#if paused}
					<PlayIcon />
				{:else}
					<PauseIcon />
				{/if}
			</button>

			<button
				type="button"
				class="rounded-full p-2 text-white transition hover:bg-white/10"
				onclick={() => seekBy(-10)}
				aria-label="Skip back 10 seconds"
			>
				<SkipPreviousIcon />
			</button>

			<button
				type="button"
				class="rounded-full p-2 text-white transition hover:bg-white/10"
				onclick={() => seekBy(10)}
				aria-label="Skip forward 10 seconds"
			>
				<SkipNextIcon />
			</button>

			<div class="ml-auto flex items-center gap-2 sm:ml-0">
				<button
					type="button"
					class="rounded-full p-2 text-white transition hover:bg-white/10"
					onclick={toggleMute}
					aria-label={muted || volume === 0 ? 'Unmute video' : 'Mute video'}
				>
					{#if muted || volume === 0}
						<MuteIcon />
					{:else}
						<VolumeIcon />
					{/if}
				</button>

				<input
					type="range"
					min="0"
					max="1"
					step="0.05"
					value={muted ? 0 : volume}
					oninput={(event) => setVolume(Number((event.currentTarget as HTMLInputElement).value))}
					class="h-1 w-24 cursor-pointer accent-white"
					aria-label="Volume"
				/>

				<button
					type="button"
					class="rounded-full p-2 text-white transition hover:bg-white/10"
					onclick={toggleFullscreen}
					aria-label={isFullscreen ? 'Exit fullscreen' : 'Enter fullscreen'}
				>
					<FullscreenIcon />
				</button>
			</div>
		</div>
	</div>
</div>
