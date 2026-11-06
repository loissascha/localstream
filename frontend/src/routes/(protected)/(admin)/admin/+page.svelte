<script lang="ts">
	import { listReencodingMovies, reencodeMovie } from '$lib/api/reencoding';
	import { auth } from '$lib/auth.svelte';
	import { movies } from '$lib/movies.svelte';
	import type { ReencodeStatus } from '$lib/types/export_types';

	let moviesWithErrors = $derived(movies.movies.filter((v) => v.had_errors));
	let moviesWithWarnings = $derived(movies.movies.filter((v) => v.had_warnings));

	let reencodingStatus = $state<ReencodeStatus[]>([]);

	async function reencodeMovieAction(movieID: string) {
		if (!auth.token) return;
		const status = await reencodeMovie(auth.token, movieID);
		console.log('status ', status);
		await reloadRunningReencodings();
	}

	async function reloadRunningReencodings() {
		if (!auth.token) return;
		const res = await listReencodingMovies(auth.token);
		reencodingStatus = res.list;
	}

	$effect(() => {
		if (!auth.initialized || !auth.token) return;
		reloadRunningReencodings();
	});
</script>

<div>
	Dashboard Page<br />
	Version: pre-alpha
</div>

{#if moviesWithErrors.length > 0}
	<div class="rounded border border-red-500">
		<div class="rounded-t bg-red-500 p-2">Movies with errors</div>
		<div class="flex flex-col gap-2 p-2">
			{#each moviesWithErrors as movie (movie.id)}
				<div class="flex justify-between">
					<div>
						{movie.name}
					</div>
					<div>
						<button
							onclick={() => {
								reencodeMovieAction(movie.id);
							}}>Reencode</button
						>
					</div>
				</div>
			{/each}
		</div>
	</div>
{/if}

{#if moviesWithWarnings.length > 0}
	<div class="rounded border border-yellow-500">
		<div class="rounded-t bg-yellow-500 p-2">Movies with warnings</div>
		<div class="flex flex-col gap-2 p-2">
			{#each moviesWithWarnings as movie (movie.id)}
				<div class="flex justify-between">
					<div>
						{movie.name}
					</div>
					<div>
						<button
							onclick={() => {
								reencodeMovieAction(movie.id);
							}}>Reencode</button
						>
					</div>
				</div>
			{/each}
		</div>
	</div>
{/if}
