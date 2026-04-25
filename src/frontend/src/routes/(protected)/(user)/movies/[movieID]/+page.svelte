<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { addMovieToCollection } from '$lib/api/collections';
	import { auth } from '$lib/auth.svelte';
	import SelectCollectionOverlay from '$lib/components/overlays/SelectCollectionOverlay.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import PlusIcon from '$lib/icons/PlusIcon.svelte';
	import { movies } from '$lib/movies.svelte';
	import DOMPurify from 'dompurify';

	const movieId = $derived(page.params.movieID ?? '');
	let movie = $derived.by(() => {
		return movies.movies.find((movie) => movie.id === movieId);
	});
	let errorMessage = $state('');

	let showAddToCollection = $state(false);
</script>

<section>
	{#if errorMessage}
		<p class="text-red-700">{errorMessage}</p>
	{/if}

	{#if movie == null}
		<span>Loading...</span>
	{:else}
		<div class="flex flex-col gap-2 md:flex-row">
			<div class="shrink-0">
				{#if movie.medium_image_url != null}
					<div>
						<img alt={movie.name} class="max-h-102" src={movie.medium_image_url} />
					</div>
				{/if}
			</div>
			<div>
				<h1 class="mb-2 px-2 text-3xl">
					{movie.name}
					{#if movie.year > 0}
						({movie.year})
					{/if}
				</h1>
				<div class="px-2">
					{@html DOMPurify.sanitize(movie.description)}
				</div>
				<div class="p-4">
					<button
						onclick={() => {
							goto(resolve('/(protected)/watch/movies/[movieID]', { movieID: movieId }));
						}}
						class="flex cursor-pointer gap-1 rounded bg-brand/80 px-4 py-2 font-semibold hover:bg-brand"
					>
						Watch <ChevronRightIcon />
					</button>
					<button
						onclick={() => {
							showAddToCollection = true;
						}}
						class="mt-4 flex cursor-pointer gap-1 rounded bg-neutral-800 px-4 py-2 font-semibold hover:bg-neutral-700"
					>
						<PlusIcon />
						Add to Collection
					</button>
				</div>
			</div>
		</div>
	{/if}
</section>

{#if showAddToCollection}
	<SelectCollectionOverlay
		selectedCollection={(collectionId) => {
			if (auth.token && movie) {
				addMovieToCollection(auth.token, collectionId, movie.id)
					.then(() => {
						showAddToCollection = false;
					})
					.catch((e) => {
						alert((e as Error).message);
					});
			}
		}}
		close={() => {
			showAddToCollection = false;
		}}
	/>
{/if}
