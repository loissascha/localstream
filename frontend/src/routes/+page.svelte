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

<main
	class="min-h-dvh bg-[radial-gradient(circle_at_15%_10%,#d7e8ff_0%,transparent_40%),radial-gradient(circle_at_85%_0%,#c8f5e9_0%,transparent_32%),linear-gradient(180deg,#eef2f7_0%,#dce6f2_100%)] px-5 py-5 text-slate-900"
>
	<header>
		<h1 class="m-0 text-[clamp(1.6rem,2.8vw,2.4rem)] font-semibold">Localstream</h1>
		<p class="mb-5 mt-1.5 text-slate-700">Your local MP4 library</p>
	</header>

	{#if loading}
		<p>Loading video library...</p>
	{:else if errorMessage}
		<p class="text-red-700">{errorMessage}</p>
	{:else if videos.length === 0}
		<p>No streamable MP4 files found in your backend video directory.</p>
	{:else}
		<section class="grid grid-cols-[repeat(auto-fill,minmax(14rem,1fr))] gap-3.5">
			{#each videos as video}
				<a
					class="grid grid-rows-[8.5rem_auto] overflow-hidden rounded-xl border border-slate-900/10 bg-white/80 no-underline transition-transform duration-150 ease-out hover:-translate-y-0.5"
					href={`/watch/${encodeURIComponent(video.id)}`}
				>
					<div
						class="flex items-center justify-center bg-[linear-gradient(135deg,#0f172a_0%,#0b4c6a_100%)] text-xl font-bold tracking-[0.08em] text-slate-50"
					>
						MP4
					</div>
					<div class="px-3.5 pb-3.5 pt-3">
						<h2 class="m-0 break-words text-[0.96rem] leading-[1.35] text-slate-900">{video.name}</h2>
						<p class="mt-1.5 text-[0.86rem] text-slate-600">{toHumanSize(video.size)}</p>
					</div>
				</a>
			{/each}
		</section>
	{/if}
</main>
