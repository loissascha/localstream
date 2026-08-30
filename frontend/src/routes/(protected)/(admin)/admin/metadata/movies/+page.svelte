<script lang="ts">
	import AdminMovieMetadataBlock from '$lib/components/admin/AdminMovieMetadataBlock.svelte';
	import { movies } from '$lib/movies.svelte';

	var errorMessage = $state('');
	var hideSingle = $state(true);
</script>

<div class="mb-4">
	<input id="hidesingle" bind:checked={hideSingle} class="cursor-pointer" type="checkbox" />
	<label for="hidesingle" class="cursor-pointer"> Hide items with single metadata</label>
</div>

{#if errorMessage != ''}
	<p class="text-red-500">{errorMessage}</p>
{/if}

<section class="mt-4 grid grid-cols-1 gap-4">
	{#each movies.movies as movie (movie.id)}
		{#if !hideSingle || movie.fetch_source == 'multiple' || movie.fetch_source == 'empty'}
			<AdminMovieMetadataBlock {movie} />
		{/if}
	{/each}
</section>
