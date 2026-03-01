<script lang="ts">
	import { page } from '$app/state';

	const videoId = $derived(page.params.id ?? '');
	const streamUrl = $derived(`/api/videos/stream?id=${encodeURIComponent(videoId)}`);
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
			class="h-auto max-h-[calc(100dvh-4.2rem)] w-full rounded-xl bg-black md:h-[min(88dvh,calc(100dvh-4.2rem))] md:w-[min(100%,120rem)]"
			controls
			autoplay
			preload="metadata"
			src={streamUrl}
		></video>
	</section>
</main>
