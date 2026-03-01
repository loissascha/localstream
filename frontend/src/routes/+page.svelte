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
	let selectedVideoId = $state('');
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

	const streamUrl = $derived(selectedVideoId ? `/api/videos/stream?id=${encodeURIComponent(selectedVideoId)}` : '');

	onMount(async () => {
		try {
			const res = await fetch('/api/videos');
			if (!res.ok) {
				throw new Error(`Failed to load videos: ${res.status}`);
			}

			const data = (await res.json()) as VideoListResponse;
			videos = data.videos ?? [];
			if (videos.length > 0) {
				selectedVideoId = videos[0].id;
			}
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Unknown error while loading videos';
		} finally {
			loading = false;
		}
	});
</script>

<main class="page">
	<section class="panel library">
		<h1>Localstream</h1>
		<p>Milestone 01: MP4 streaming via HTTP range requests</p>

		{#if loading}
			<p>Loading video library...</p>
		{:else if errorMessage}
			<p class="error">{errorMessage}</p>
		{:else if videos.length === 0}
			<p>No streamable MP4 files found in your backend video directory.</p>
		{:else}
			<ul>
				{#each videos as video}
					<li>
						<button
							type="button"
							class:selected={selectedVideoId === video.id}
							onclick={() => {
								selectedVideoId = video.id;
							}}
						>
							<span>{video.name}</span>
							<small>{toHumanSize(video.size)}</small>
						</button>
					</li>
				{/each}
			</ul>
		{/if}
	</section>

	<section class="panel player">
		{#if streamUrl}
			<!-- svelte-ignore a11y_media_has_caption -->
			<video controls preload="metadata" src={streamUrl}></video>
		{:else}
			<p>Select a video from the library to start playback.</p>
		{/if}
	</section>
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: 'IBM Plex Sans', 'Segoe UI', sans-serif;
		background: linear-gradient(180deg, #eef2f7 0%, #dce6f2 100%);
		color: #0f172a;
	}

	.page {
		display: grid;
		grid-template-columns: minmax(18rem, 24rem) 1fr;
		gap: 1rem;
		padding: 1rem;
		min-height: 100dvh;
		box-sizing: border-box;
	}

	.panel {
		background: rgba(255, 255, 255, 0.86);
		border: 1px solid rgba(15, 23, 42, 0.08);
		border-radius: 0.75rem;
		padding: 1rem;
		backdrop-filter: blur(4px);
	}

	.library h1 {
		margin: 0;
		font-size: 1.5rem;
	}

	.library p {
		margin: 0.5rem 0 1rem;
		color: #334155;
	}

	ul {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	button {
		width: 100%;
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.75rem;
		padding: 0.65rem 0.75rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.5rem;
		background: #f8fafc;
		color: #0f172a;
		text-align: left;
		cursor: pointer;
	}

	button:hover {
		background: #f1f5f9;
	}

	button.selected {
		border-color: #0369a1;
		background: #e0f2fe;
	}

	small {
		color: #475569;
		white-space: nowrap;
	}

	.player {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	video {
		width: min(100%, 70rem);
		max-height: min(80dvh, 42rem);
		background: #020617;
		border-radius: 0.75rem;
	}

	.error {
		color: #b91c1c;
	}

	@media (max-width: 900px) {
		.page {
			grid-template-columns: 1fr;
		}

		video {
			max-height: 55dvh;
		}
	}
</style>
