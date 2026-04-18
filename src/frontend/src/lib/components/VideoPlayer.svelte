<script lang="ts">
	import { onDestroy, type Snippet } from 'svelte';
	import FullscreenIcon from '$lib/icons/FullscreenIcon.svelte';
	import MuteIcon from '$lib/icons/MuteIcon.svelte';
	import PauseIcon from '$lib/icons/PauseIcon.svelte';
	import PlayIcon from '$lib/icons/PlayIcon.svelte';
	import SkipNextIcon from '$lib/icons/SkipNextIcon.svelte';
	import SkipPreviousIcon from '$lib/icons/SkipPreviousIcon.svelte';
	import VolumeIcon from '$lib/icons/VolumeIcon.svelte';
	import { getCookie, setCookie } from '$lib/cookies';
	import FullscreenExitIcon from '$lib/icons/FullscreenExitIcon.svelte';

	interface OverlayState {
		currentTime: number;
		duration: number;
		paused: boolean;
		isFullscreen: boolean;
	}

	interface Props {
		href: string;
		duration?: number;
		currentTime?: number;
		onplay?: () => void;
		onpause?: () => void;
		onended?: () => void;
		overlay?: Snippet<[OverlayState]>;
		topbar?: Snippet;
	}

	let {
		href,
		onplay,
		onpause,
		onended,
		overlay,
		topbar,
		duration = $bindable(0),
		currentTime = $bindable(0)
	}: Props = $props();

	let containerEl = $state<HTMLDivElement | null>(null);
	let videoEl = $state<HTMLVideoElement | null>(null);
	let paused = $state(true);
	let muted = $state(false);
	let volume = $state(1);
	let isFullscreen = $state(false);
	let seekValue = $state(0);
	let showControls = $state(true);
	let hideControlsTimer: ReturnType<typeof setTimeout> | null = null;

	function syncState() {
		if (!videoEl) return;
		currentTime = videoEl.currentTime;
		seekValue = videoEl.currentTime;
		duration = Number.isFinite(videoEl.duration) ? videoEl.duration : 0;
		paused = videoEl.paused;
		muted = videoEl.muted;
		volume = videoEl.volume;
	}

	function clearHideControlsTimer() {
		if (hideControlsTimer !== null) {
			clearTimeout(hideControlsTimer);
			hideControlsTimer = null;
		}
	}

	function scheduleHideControls() {
		clearHideControlsTimer();
		if (paused) return;
		hideControlsTimer = setTimeout(() => {
			showControls = false;
			hideControlsTimer = null;
		}, 2000);
	}

	function revealControls() {
		showControls = true;
		scheduleHideControls();
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
		revealControls();
		if (!videoEl) return;
		if (videoEl.paused) {
			await videoEl.play();
			return;
		}
		videoEl.pause();
	}

	function seekTo(value: number) {
		revealControls();
		if (!videoEl) return;
		const boundedValue = Math.min(Math.max(value, 0), duration || 0);
		videoEl.currentTime = boundedValue;
		currentTime = boundedValue;
		seekValue = boundedValue;
		syncState();
	}

	const seekMax = $derived(Math.max(duration, currentTime, seekValue, 0));

	function seekBy(seconds: number) {
		seekTo(currentTime + seconds);
	}

	function setVolume(value: number) {
		revealControls();
		if (!videoEl) return;
		const boundedValue = Math.min(Math.max(value, 0), 1);
		videoEl.volume = boundedValue;
		videoEl.muted = boundedValue === 0;
		volume = boundedValue;
		muted = videoEl.muted;
		setCookie('videoplayer_volume', boundedValue.toString(), 300);
	}

	function toggleMute() {
		revealControls();
		if (!videoEl) return;
		videoEl.muted = !videoEl.muted;
		muted = videoEl.muted;
	}

	async function toggleFullscreen() {
		revealControls();
		if (!containerEl) return;
		if (document.fullscreenElement === containerEl) {
			await document.exitFullscreen();
			return;
		}
		await containerEl.requestFullscreen();
	}

	function handlePlay() {
		syncState();
		revealControls();
		onplay?.();
	}

	function handlePause() {
		syncState();
		showControls = true;
		clearHideControlsTimer();
		onpause?.();
	}

	function handleEnded() {
		syncState();
		showControls = true;
		clearHideControlsTimer();
		onended?.();
	}

	function handleFullscreenChange() {
		isFullscreen = document.fullscreenElement === containerEl;
		revealControls();
	}

	function handleInteraction() {
		revealControls();
	}

	function isInteractiveTarget(target: EventTarget | null) {
		return target instanceof HTMLElement && target.closest('button, input, a') !== null;
	}

	function focusContainer() {
		containerEl?.focus({ preventScroll: true });
	}

	function handlePointerDown(event: PointerEvent) {
		handleInteraction();
		if (!isInteractiveTarget(event.target)) {
			focusContainer();
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (isInteractiveTarget(event.target)) {
			return;
		}

		switch (event.key) {
			case ' ':
			case 'k':
			case 'K':
				event.preventDefault();
				void togglePlay();
				break;
			case 'ArrowLeft':
				if (event.shiftKey) return;
				event.preventDefault();
				seekBy(-10);
				break;
			case 'ArrowRight':
				if (event.shiftKey) return;
				event.preventDefault();
				seekBy(10);
				break;
			case 'ArrowUp':
				event.preventDefault();
				setVolume((muted ? 0 : volume) + 0.05);
				break;
			case 'ArrowDown':
				event.preventDefault();
				setVolume((muted ? 0 : volume) - 0.05);
				break;
			case 'm':
			case 'M':
				event.preventDefault();
				toggleMute();
				break;
			case 'f':
			case 'F':
				event.preventDefault();
				void toggleFullscreen();
				break;
		}
	}

	$effect(() => {
		if (videoEl && Math.abs(videoEl.currentTime - currentTime) > 0.25) {
			videoEl.currentTime = currentTime;
		}
		seekValue = currentTime;
	});

	$effect(() => {
		document.addEventListener('fullscreenchange', handleFullscreenChange);

		return () => {
			document.removeEventListener('fullscreenchange', handleFullscreenChange);
		};
	});

	$effect(() => {
		const volumeCookie = getCookie('videoplayer_volume');
		if (volumeCookie != null && !isNaN(Number(volumeCookie))) {
			setVolume(Number(volumeCookie));
		}
	});

	onDestroy(() => {
		clearHideControlsTimer();
		if (document.fullscreenElement === containerEl) {
			document.exitFullscreen().catch(() => {});
		}
	});
</script>

<!-- svelte-ignore a11y_media_has_caption -->
<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<div
	bind:this={containerEl}
	tabindex="0"
	role="application"
	aria-label="Video player"
	class={`relative h-full w-full overflow-hidden bg-black outline-none ${showControls || paused ? 'cursor-default' : 'cursor-none'}`}
	onmousemove={handleInteraction}
	onpointerdown={handlePointerDown}
	ontouchstart={handleInteraction}
	onkeydown={handleKeydown}
>
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
		onseeking={syncState}
		disablepictureinpicture
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

	{#if topbar}
		<div
			class={`absolute top-0 right-0 left-0 z-10 flex items-center gap-2 bg-linear-to-t via-black/55 to-black/90 px-4 py-4 transition-opacity duration-200 ${showControls || paused ? 'opacity-100' : 'pointer-events-none opacity-0'}`}
		>
			{@render topbar()}
		</div>
	{/if}

	{#if overlay}
		<div class={`absolute right-4 bottom-24 z-10 max-w-[min(20rem,calc(100%-2rem))]`}>
			{@render overlay({ currentTime, duration, paused, isFullscreen })}
		</div>
	{/if}

	<div
		class={`absolute inset-x-0 bottom-0 bg-linear-to-t from-black/90 via-black/55 to-transparent px-4 pt-10 pb-4 text-white transition-opacity duration-200 ${showControls || paused ? 'opacity-100' : 'pointer-events-none opacity-0'}`}
	>
		<div class="mb-3 flex items-center gap-3 text-xs text-white/80">
			<span class="w-12 text-right tabular-nums">{formatTime(currentTime)}</span>
			<input
				type="range"
				min="0"
				max={seekMax}
				step="0.1"
				bind:value={seekValue}
				oninput={(event) => seekTo(Number((event.currentTarget as HTMLInputElement).value))}
				class="h-1 w-full cursor-pointer accent-white"
				aria-label="Seek"
			/>
			<span class="w-12 tabular-nums">{formatTime(duration)}</span>
		</div>

		<div class="flex grow flex-wrap items-center justify-between gap-2 sm:flex-nowrap">
			<div>
				<button
					type="button"
					class="cursor-pointer rounded-full p-2 text-white transition hover:bg-white/10"
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
					class="cursor-pointer rounded-full p-2 text-white transition hover:bg-white/10"
					onclick={() => seekBy(-10)}
					aria-label="Skip back 10 seconds"
				>
					<SkipPreviousIcon />
				</button>

				<button
					type="button"
					class="cursor-pointer rounded-full p-2 text-white transition hover:bg-white/10"
					onclick={() => seekBy(10)}
					aria-label="Skip forward 10 seconds"
				>
					<SkipNextIcon />
				</button>
			</div>

			<div class="ml-auto flex items-center gap-2 sm:ml-0">
				<button
					type="button"
					class="cursor-pointer rounded-full p-2 text-white transition hover:bg-white/10"
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
					class="cursor-pointer rounded-full p-2 text-white transition hover:bg-white/10"
					onclick={toggleFullscreen}
					aria-label={isFullscreen ? 'Exit fullscreen' : 'Enter fullscreen'}
				>
					{#if isFullscreen}
						<FullscreenExitIcon />
					{:else}
						<FullscreenIcon />
					{/if}
				</button>
			</div>
		</div>
	</div>
</div>
