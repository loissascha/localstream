<script lang="ts">
	import {
		computePosition,
		flip,
		offset as floatingOffset,
		shift,
		type VirtualElement
	} from '@floating-ui/dom';
	import { formatTime } from '$lib/format';
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
		backlink?: string;
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
		backlink,
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

	// seek bar
	let seekBarEl = $state<HTMLDivElement | null>(null);
	let seekTooltipEl = $state<HTMLDivElement | null>(null);
	let showSeekTooltip = $state(false);
	let hoverSeekTime = $state(0);
	let hoverSeekX = $state(0);
	let hoverSeekY = $state(0);
	let seekTooltipX = $state(0);
	let seekTooltipY = $state(0);

	const seekMax = $derived(Math.max(duration, currentTime, seekValue, 0));

	$effect(() => {
		if (videoEl && Math.abs(videoEl.currentTime - currentTime) > 0.25) {
			videoEl.currentTime = currentTime;
		}
		seekValue = currentTime;
	});

	function mouseMoved() {
		console.log('mouse moved');
		revealControls();
	}

	function mouseClicked() {
		console.log('mouse clicked');
		if (showControls) {
			if (paused) {
				play();
			} else {
				pause();
			}
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
		console.log('play');
	}

	async function pause() {
		if (!videoEl) return;
		videoEl.pause();
		onpause?.();
		paused = true;
		revealControls();
		syncState();
		console.log('pause');
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
		}, 4000);
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

	function getPercentageBetween(start: number, end: number, value: number): number {
		if (end <= start) return 0;
		return ((value - start) / (end - start)) * 100;
	}

	function handleSeekPointerEnter(event: PointerEvent) {
		showSeekTooltip = true;
		updateSeekHover(event);
	}

	function handleSeekPointerMove(event: PointerEvent) {
		if (!showSeekTooltip) return;

		updateSeekHover(event);
	}

	function handleSeekPointerLeave() {
		showSeekTooltip = false;
	}

	function getSeekTimeFromPointer(event: PointerEvent) {
		if (!seekBarEl || seekMax <= 0) return 0;

		const rect = seekBarEl.getBoundingClientRect();
		if (rect.width <= 0) return 0;

		const ratio = Math.min(Math.max((event.clientX - rect.left) / rect.width, 0), 1);
		hoverSeekX = Math.min(Math.max(event.clientX, rect.left), rect.right);
		hoverSeekY = rect.top;

		return ratio * seekMax;
	}

	function updateSeekHover(event: PointerEvent) {
		hoverSeekTime = getSeekTimeFromPointer(event);
		void updateSeekTooltipPosition();
	}

	async function updateSeekTooltipPosition() {
		if (!seekTooltipEl || !showSeekTooltip) return;

		const { x, y, strategy } = await computePosition(createSeekVirtualReference(), seekTooltipEl, {
			placement: 'top',
			strategy: 'fixed',
			middleware: [floatingOffset(10), flip(), shift({ padding: 8 })]
		});

		seekTooltipX = x;
		seekTooltipY = y;
	}

	function createSeekVirtualReference(): VirtualElement {
		return {
			getBoundingClientRect() {
				return {
					width: 0,
					height: 0,
					x: hoverSeekX,
					y: hoverSeekY,
					top: hoverSeekY,
					right: hoverSeekX,
					bottom: hoverSeekY,
					left: hoverSeekX
				};
			}
		};
	}

	function handleSeekPointerDown(event: PointerEvent) {
		preventClick(event);
		pause();
		seekTo(hoverSeekTime);
	}

	function seekTo(value: number) {
		if (!videoEl) return;
		const boundedValue = Math.min(Math.max(value, 0), duration || 0);
		videoEl.currentTime = boundedValue;
		currentTime = boundedValue;
		seekValue = boundedValue;
		syncState();
	}

	function forward10Seconds() {
		seekTo(currentTime + 10);
	}

	function backdward10Seconds() {
		seekTo(currentTime - 10);
	}

	function preventClick(e: PointerEvent) {
		e.preventDefault();
		e.stopPropagation();
	}
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
			<div class="flex items-center gap-8">
				<button
					class="pointer-events-auto flex cursor-pointer items-center justify-center rounded-full bg-neutral-800/60"
					onpointerdown={preventClick}
					onclick={() => {
						backdward10Seconds();
					}}
				>
					<ChevronLeftIcon size={50} />
				</button>
				<button
					onpointerdown={preventClick}
					onclick={() => {
						if (paused) {
							play();
						} else {
							pause();
						}
					}}
					class="pointer-events-auto flex cursor-pointer items-center justify-center rounded-full bg-neutral-800/60"
				>
					{#if paused}
						<PlayIcon size={120} />
					{:else}
						<PauseIcon size={120} />
					{/if}
				</button>
				<button
					class="pointer-events-auto flex cursor-pointer items-center justify-center rounded-full bg-neutral-800/60"
					onpointerdown={preventClick}
					onclick={() => {
						forward10Seconds();
					}}
				>
					<ChevronRightIcon size={50} />
				</button>
			</div>
		</div>

		<!-- Top Bar -->
		<div class="pointer-events-none absolute top-0 right-0 left-0">
			<div class="flex items-center justify-between px-8 py-4">
				<div class="pointer-events-auto">
					{#if backlink}
						<a href={backlink}><ChevronLeftIcon size={40} /></a>
					{/if}
				</div>
				<div class="pointer-events-auto">Right</div>
			</div>
		</div>

		<!-- Bottom Bar -->
		<div class="pointer-events-none absolute right-0 bottom-0 left-0 px-4 py-2">
			<!-- first bottom bar -->
			<div class="flex items-center justify-between gap-4">
				<div class="pointer-events-auto shrink-0 text-sm">{formatTime(currentTime)}</div>
				<div class="pointer-events-auto grow">
					{#if showSeekTooltip}
						<div
							bind:this={seekTooltipEl}
							class="pointer-events-none z-20 rounded-md bg-neutral-700/85 px-2 py-1 text-xs font-medium text-white tabular-nums shadow-lg ring-1 ring-white/10 backdrop-blur-sm"
							style={`position: fixed; left: ${seekTooltipX}px; top: ${seekTooltipY}px;`}
						>
							{formatTime(hoverSeekTime)}
						</div>
					{/if}

					<!-- Seek Bar -->
					<!-- svelte-ignore a11y_no_static_element_interactions -->
					<div
						bind:this={seekBarEl}
						onpointerenter={handleSeekPointerEnter}
						onpointermove={handleSeekPointerMove}
						onpointerleave={handleSeekPointerLeave}
						onpointerdown={handleSeekPointerDown}
						class="group relative h-2 w-full cursor-pointer rounded-full bg-neutral-600"
					>
						<div
							class="absolute h-2 rounded-full bg-neutral-500"
							style={`width: ${getPercentageBetween(0, seekMax, bufferedUntil)}%;`}
						></div>
						<div
							class="absolute h-2 rounded-full bg-brand"
							style={`width: ${getPercentageBetween(0, seekMax, seekValue)}%;`}
						></div>
						<div
							class="absolute top-1 h-5 w-5 -translate-y-1/2 rounded-full border-2 border-brand bg-neutral-500 transition-all duration-300 group-hover:h-6 group-hover:w-6"
							style={`left: calc(${getPercentageBetween(0, seekMax, seekValue)}% - 10px);`}
						></div>
					</div>
				</div>
				<div class="pointer-events-auto shrink-0 text-sm">{formatTime(duration)}</div>
			</div>
			<!-- second bottom bar -->
			<div class="flex items-center justify-between">
				<div class="pointer-events-auto">
					<button
						onpointerdown={preventClick}
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
							<PlayIcon size={30} />
						{:else}
							<PauseIcon size={30} />
						{/if}
					</button>
				</div>
				<div class="pointer-events-auto">Right</div>
			</div>
		</div>
	</div>
</div>
