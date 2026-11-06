<script lang="ts">
	import { listReencodingMovies, reencodeMovie } from '$lib/api/reencoding';
	import { auth } from '$lib/auth.svelte';
	import { movies } from '$lib/movies.svelte';
	import type { MovieInfo, ReencodeStatus } from '$lib/types/export_types';

	interface MovieRes {
		mi: MovieInfo;
		rs?: ReencodeStatus;
	}

	let moviesWithErrors = $derived(movies.movies.filter((v) => v.had_errors));
	let moviesWithWarnings = $derived(movies.movies.filter((v) => v.had_warnings));
	let reencodingStatus = $state<ReencodeStatus[]>([]);

	let moviesWithErrorsWithStatus = $state<MovieRes[]>([]);
	let moviesWithWarningsWithStatus = $state<MovieRes[]>([]);

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

	$effect(() => {
		let ms: MovieRes[] = [];
		for (let movie of moviesWithErrors) {
			let rs = reencodingStatus.find((f) => f.id == movie.id);
			if (rs) {
				ms.push({
					mi: movie,
					rs: rs
				});
			} else {
				ms.push({
					mi: movie
				});
			}
		}
		moviesWithErrorsWithStatus = ms;
		ms = [];
		for (let movie of moviesWithWarnings) {
			let rs = reencodingStatus.find((f) => f.id == movie.id);
			if (rs) {
				ms.push({
					mi: movie,
					rs: rs
				});
			} else {
				ms.push({
					mi: movie
				});
			}
		}
		moviesWithWarningsWithStatus = ms;
	});
</script>

<div>
	Dashboard Page<br />
	Version: pre-alpha
</div>

{#if moviesWithErrorsWithStatus.length > 0}
	<div class="mt-4 rounded border border-red-500">
		<div class="rounded-t bg-red-500 p-2">Movies with errors</div>
		<div class="flex flex-col">
			{#each moviesWithErrorsWithStatus as movie (movie.mi.id)}
				<div class="flex justify-between p-2 odd:bg-neutral-800">
					<div>
						{movie.mi.name}
					</div>
					<div>
						{#if movie.rs}
							<div>
								Started: {movie.rs.started} Finished {movie.rs.finished} Failed {movie.rs.failed} Error
								{movie.rs.error}
							</div>
						{:else}
							<button
								onclick={() => {
									reencodeMovieAction(movie.mi.id);
								}}>Reencode</button
							>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	</div>
{/if}

{#if moviesWithWarnings.length > 0}
	<div class="mt-4 rounded border border-yellow-500">
		<div class="rounded-t bg-yellow-500 p-2">Movies with warnings</div>
		<div class="flex flex-col gap-2 p-2">
			{#each moviesWithWarnings as movie (movie.id)}
				<div class="flex justify-between p-2 odd:bg-neutral-800">
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
