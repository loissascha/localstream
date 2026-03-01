<script lang="ts">
	import { onMount } from 'svelte';

	type VideoListItem = {
		id: string;
		name: string;
		size: number;
		mimeType: string;
	};

	type VideoListResponse = {
		videos: VideoListItem[];
	};

	let videos = $state<VideoListItem[]>([]);
	let loading = $state(true);
	let errorMessage = $state('');

	const toHumanSize = (bytes: number): string => {
		if (bytes < 1024) {
			return `${bytes} B`;
		}

		const units = ['KB', 'MB', 'GB', 'TB'];
		let size = bytes / 1024;
		let unitIndex = 0;

		while (size >= 1024 && unitIndex < units.length - 1) {
			size /= 1024;
			unitIndex += 1;
		}

		return `${size.toFixed(1)} ${units[unitIndex]}`;
	};

	onMount(async () => {
		try {
			const res = await fetch('/api/videos');
			if (!res.ok) {
				throw new Error(`Failed to load videos: ${res.status}`);
			}

			const data = (await res.json()) as VideoListResponse;
			videos = data.videos ?? [];
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loading = false;
		}
	});
</script>

<main class="page">
	<header>
		<h1>Localstream</h1>
		<p>Your local MP4 library</p>
	</header>

	{#if loading}
		<p>Loading video library...</p>
	{:else if errorMessage}
		<p class="error">{errorMessage}</p>
	{:else if videos.length === 0}
		<p>No streamable MP4 files found in your backend video directory.</p>
	{:else}
		<section class="grid">
			{#each videos as video}
				<a class="card" href={`/watch/${encodeURIComponent(video.id)}`}>
					<div class="thumb">MP4</div>
					<div class="meta">
						<h2>{video.name}</h2>
						<p>{toHumanSize(video.size)}</p>
					</div>
				</a>
			{/each}
		</section>
	{/if}
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: 'IBM Plex Sans', 'Segoe UI', sans-serif;
		background:
			radial-gradient(circle at 15% 10%, #d7e8ff 0%, transparent 40%),
			radial-gradient(circle at 85% 0%, #c8f5e9 0%, transparent 32%),
			linear-gradient(180deg, #eef2f7 0%, #dce6f2 100%);
		color: #0f172a;
	}

	.page {
		padding: 1.25rem;
		min-height: 100dvh;
		box-sizing: border-box;
	}

	header h1 {
		margin: 0;
		font-size: clamp(1.6rem, 2.8vw, 2.4rem);
	}

	header p {
		margin: 0.35rem 0 1.25rem;
		color: #334155;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(14rem, 1fr));
		gap: 0.9rem;
	}

	.card {
		display: grid;
		grid-template-rows: 8.5rem auto;
		text-decoration: none;
		background: rgba(255, 255, 255, 0.82);
		border: 1px solid rgba(15, 23, 42, 0.08);
		border-radius: 0.8rem;
		overflow: hidden;
		transition: transform 140ms ease;
	}

	.card:hover {
		transform: translateY(-2px);
	}

	.thumb {
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.2rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		color: #f8fafc;
		background: linear-gradient(135deg, #0f172a 0%, #0b4c6a 100%);
	}

	.meta {
		padding: 0.7rem 0.8rem 0.85rem;
	}

	.meta h2 {
		margin: 0;
		font-size: 0.96rem;
		color: #0f172a;
		line-height: 1.35;
		word-break: break-word;
	}

	.meta p {
		margin: 0.35rem 0 0;
		font-size: 0.86rem;
		color: #475569;
	}

	.error {
		color: #b91c1c;
	}
</style>
