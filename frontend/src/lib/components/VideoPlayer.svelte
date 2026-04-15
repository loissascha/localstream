<script lang="ts">
	interface Props {
		href: string;
		onplay: () => void;
		onpause: () => void;
		onended: () => void;
	}
	let { href, onplay, onpause, onended }: Props = $props();

	let videoEl = $state<HTMLVideoElement | null>(null);
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
	></video>
</div>
