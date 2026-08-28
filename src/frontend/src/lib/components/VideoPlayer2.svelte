<script lang="ts">
	import ChevronLeftIcon from '$lib/icons/ChevronLeftIcon.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import PauseIcon from '$lib/icons/PauseIcon.svelte';
	import PlayIcon from '$lib/icons/PlayIcon.svelte';
	import type { SubtitleInfo } from '$lib/types/export_types';
	import { onDestroy, type Snippet } from 'svelte';

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
		bottomrightextensions?: Snippet;
		subtitles?: SubtitleInfo[];
	}

	let {
		href,
		onplay,
		onpause,
		onended,
		overlay,
		topbar,
		bottomrightextensions,
		subtitles,
		duration = $bindable(0),
		currentTime = $bindable(0)
	}: Props = $props();

	let showControls = $state(true);
	let paused = $state(true);
	let videoEl = $state<HTMLVideoElement | null>(null);
	let hideControlsTimer: ReturnType<typeof setTimeout> | null = null;
	let muted = $state(false);
	let volume = $state(1);
	let seekValue = $state(0);
	let bufferedUntil = $state(0);

	function mouseMoved() {
		console.log('mouse moved');
		revealControls();
	}

	function mouseClicked() {
		console.log('mouse clicked');
		if (showControls) {
			// if (paused) {
			// 	play();
			// } else {
			// 	pause();
			// }
		} else {
			revealControls();
		}
	}

	async function play() {
		if (!videoEl) return;
		await videoEl.play();
		onplay?.();
		paused = false;
		scheduleHideControls();
		syncState();
	}

	async function pause() {
		if (!videoEl) return;
		videoEl.pause();
		onpause?.();
		paused = true;
		revealControls();
		syncState();
	}

	function revealControls() {
		showControls = true;
		scheduleHideControls();
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

	function syncState() {
		if (!videoEl) return;
		currentTime = videoEl.currentTime;
		seekValue = videoEl.currentTime;
		duration = Number.isFinite(videoEl.duration) ? videoEl.duration : 0;
		syncBuffered();
	}

	function syncBuffered() {
		if (!videoEl || videoEl.buffered.length === 0) {
			bufferedUntil = 0;
			return;
		}
		bufferedUntil = videoEl.buffered.end(videoEl.buffered.length - 1);
	}

	onDestroy(() => {
		clearHideControlsTimer();
	});
</script>

<!-- svelte-ignore a11y_media_has_caption -->
<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<div
	tabindex="0"
	role="application"
	aria-label="Video player"
	class={`relative h-full w-full overflow-hidden bg-black outline-none ${showControls || paused ? 'cursor-default' : 'cursor-none'}`}
	onmousemove={mouseMoved}
	onpointerdown={mouseClicked}
>
	<video
		bind:this={videoEl}
		class="h-full w-full bg-black object-contain"
		preload="metadata"
		playsinline
		src={href}
		onloadedmetadata={syncState}
		ondurationchange={syncState}
		onratechange={syncState}
		onseeked={syncState}
		onvolumechange={syncState}
		ontimeupdate={syncState}
		onseeking={syncState}
		disablepictureinpicture
	>
		{#each subtitles ?? [] as subtitle}
			<track
				src={subtitle.path}
				kind="subtitles"
				srclang={subtitle.lang_short}
				label={subtitle.lang}
			/>
		{/each}
	</video>

	<!-- Controls -->
	<div
		class={`absolute top-0 right-0 bottom-0 left-0 ${showControls ? 'opacity-100' : 'pointer-events-none opacity-0'}`}
	>
		<!-- Display Center Buttons -->
		<div class="pointer-events-none absolute inset-0 flex items-center justify-center">
			<div class="flex items-center gap-1">
				<button>
					<ChevronLeftIcon size={50} />
				</button>
				<button
					onclick={() => {
						if (paused) {
							play();
						} else {
							pause();
						}
					}}
					class="pointer-events-auto cursor-pointer"
				>
					{#if paused}
						<PlayIcon size={120} />
					{:else}
						<PauseIcon size={120} />
					{/if}
				</button>
				<button>
					<ChevronRightIcon size={50} />
				</button>
			</div>
		</div>

		<!-- Top Bar -->
		<div class="pointer-events-none absolute top-0 right-0 left-0">
			<div class="flex items-center justify-between">
				<div class="pointer-events-auto">Left</div>
				<div class="pointer-events-auto">Right</div>
			</div>
		</div>

		<!-- Bottom Bar -->
		<div class="pointer-events-none absolute right-0 bottom-0 left-0">
			<div class="flex items-center justify-between">
				<div class="pointer-events-auto">Left</div>
				<div class="pointer-events-auto">Center</div>
				<div class="pointer-events-auto">Right</div>
			</div>
		</div>
	</div>
</div>
