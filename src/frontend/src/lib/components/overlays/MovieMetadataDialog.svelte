<script lang="ts">
	import type { MovieInfo, MovieResult } from '$lib/types/export_types';
	import type { Snippet } from 'svelte';
	import Dialog from './Dialog.svelte';
	import { auth } from '$lib/auth.svelte';
	import { searchMovieMetadata, setPrimaryMovieMetadataByFetchID } from '$lib/api/movie_metadata';
	import SearchIcon from '$lib/icons/SearchIcon.svelte';

	interface Props {
		movie: MovieInfo;
		children: Snippet;
		updated: () => void;
	}
	let { movie, children, updated }: Props = $props();

	let searchingMetadata = $state(false);
	let error_message = $state('');
	let success_message = $state('');
	let searchQuery = $state('');
	let searchResults = $state<MovieResult[]>([]);

	async function submitMetadataSearchForm() {
		try {
			if (!auth.token) return;
			searchingMetadata = true;
			success_message = '';
			error_message = '';
			if (searchQuery == '') {
				alert('No query');
				return;
			}
			searchResults = await searchMovieMetadata(auth.token, searchQuery);
			searchingMetadata = false;
		} catch (e) {
			const m = (e as Error).message;
			error_message = m;
		}
	}
</script>

<Dialog>
	{@render children()}
	{#snippet content(close)}
		<h1 class="text-2xl font-bold tracking-wide">Metadata</h1>
		<h2 class="text-lg font-semibold tracking-wide">{movie.name}</h2>
		<!-- <div class="flex"> -->
		<!-- 	<div class="grow"> -->
		<!-- 		{movie.description} -->
		<!-- 	</div> -->
		<!-- 	<div> -->
		<!-- 		<img src={movie.medium_image_url} class="w-120" alt={movie.name} /> -->
		<!-- 	</div> -->
		<!-- </div> -->
		{#if success_message != ''}
			<div class="mt-4 text-green-500">
				{success_message}
			</div>
		{/if}
		{#if error_message != ''}
			<div class="mt-4 text-red-500">
				{error_message}
			</div>
		{/if}

		<form
			onsubmit={(e) => {
				e.preventDefault();
				console.log('submit form');
				if (searchingMetadata) {
					alert('Searching...');
					return;
				}
				submitMetadataSearchForm();
			}}
			class="flex items-center gap-2"
		>
			<input
				bind:value={searchQuery}
				type="text"
				class="my-4 w-full rounded bg-neutral-700 px-4 py-2"
				placeholder="Search for new Metadata"
			/>
			<button
				class="flex cursor-pointer gap-2 rounded bg-neutral-700 px-4 py-2 hover:bg-neutral-600"
			>
				<SearchIcon /> Search
			</button>
		</form>

		{#if searchingMetadata}
			Searching...
		{:else}
			{#each searchResults as result (result.id)}
				<div class="my-4 overflow-y-auto p-4">
					<div class="grid grid-cols-1 md:grid-cols-2">
						<div>
							<div class="font-bold">{result.original_title} ({result.release_year})</div>
							<div class="my-2">
								{result.overview}
							</div>
							<div>
								<button
									onclick={() => {
										if (!auth.token) return;
										setPrimaryMovieMetadataByFetchID(auth.token, movie.id, result.id)
											.then(() => {
												searchResults = [];
												searchQuery = '';
												success_message = 'Updated primary metadata!';
												updated();
												close();
											})
											.catch((e) => {
												const m = (e as Error).message;
												alert(m);
											});
									}}
									class="mt-4 mb-4 cursor-pointer rounded bg-neutral-700 px-4 py-2 hover:bg-neutral-600"
									>Select as Primary</button
								>
							</div>
						</div>
						<div>
							{#if result.poster_path != ''}
								<img
									alt={'movie ' + result.id}
									src={`https://image.tmdb.org/t/p/w500/${result.poster_path}`}
								/>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		{/if}
	{/snippet}
</Dialog>
