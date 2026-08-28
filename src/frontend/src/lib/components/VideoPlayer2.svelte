<script lang="ts">
	import type { SubtitleInfo } from '$lib/types/export_types';
	import type { Snippet } from 'svelte';

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

	let showControls = $state(false);
	let paused = $state(true);
	let videoEl = $state<HTMLVideoElement | null>(null);

	function mouseMoved() {}

	function mouseClicked() {}
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
</div>
