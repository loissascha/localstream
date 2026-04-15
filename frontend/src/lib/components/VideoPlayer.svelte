<script lang="ts">
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

	let videoEl = $state<HTMLVideoElement | null>(null);

	function syncState() {
		if (!videoEl) return;
		currentTime = videoEl.currentTime;
		duration = Number.isFinite(videoEl.duration) ? videoEl.duration : 0;
	}

	$effect(() => {
		if (videoEl && Math.abs(videoEl.currentTime - currentTime) > 0.25) {
			videoEl.currentTime = currentTime;
		}
	});
</script>

<!-- svelte-ignore a11y_media_has_caption -->
<div class="h-full w-full">
	<video
		bind:this={videoEl}
		class="h-full w-full bg-black object-contain"
		controls
		preload="metadata"
		src={href}
		{onplay}
		{onpause}
		{onended}
		ontimeupdate={syncState}
		onloadedmetadata={syncState}
	></video>
</div>
