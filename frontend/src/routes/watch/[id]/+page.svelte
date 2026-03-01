<script lang="ts">
	import { page } from '$app/state';

	const videoId = $derived(page.params.id ?? '');
	const streamUrl = $derived(`/api/videos/stream?id=${encodeURIComponent(videoId)}`);
</script>

<main class="player-page">
	<header class="topbar">
		<a class="back" href="/">Back to library</a>
	</header>

	<section class="player-wrap">
		<!-- svelte-ignore a11y_media_has_caption -->
		<video controls autoplay preload="metadata" src={streamUrl}></video>
	</section>
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: 'IBM Plex Sans', 'Segoe UI', sans-serif;
		background: #020617;
		color: #f8fafc;
	}

	.player-page {
		min-height: 100dvh;
		display: grid;
		grid-template-rows: auto 1fr;
	}

	.topbar {
		padding: 0.85rem 1rem;
		background: linear-gradient(180deg, rgba(15, 23, 42, 0.95), rgba(15, 23, 42, 0.72));
		border-bottom: 1px solid rgba(148, 163, 184, 0.18);
	}

	.back {
		display: inline-block;
		padding: 0.4rem 0.6rem;
		border-radius: 0.45rem;
		color: #cbd5e1;
		text-decoration: none;
		border: 1px solid rgba(148, 163, 184, 0.28);
	}

	.back:hover {
		color: #f8fafc;
		border-color: rgba(203, 213, 225, 0.62);
	}

	.player-wrap {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: clamp(0.5rem, 1.2vw, 1rem);
	}

	video {
		width: min(100%, 120rem);
		height: min(88dvh, calc(100dvh - 4.2rem));
		background: #000;
		border-radius: 0.75rem;
	}

	@media (max-width: 700px) {
		video {
			height: auto;
			max-height: calc(100dvh - 4.2rem);
		}
	}
</style>
