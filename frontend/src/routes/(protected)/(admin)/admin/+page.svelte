<script lang="ts">
	import { movies } from '$lib/movies.svelte';

	let moviesWithErrors = $derived(movies.movies.filter((v) => v.had_errors));
	let moviesWithWarnings = $derived(movies.movies.filter((v) => v.had_warnings));

	// $effect(() => {
	// 	console.log("there are " + movies.movies.length + " movies in existenze")
	// })
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
				<div>
					{movie.name}
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
				<div>
					{movie.name}
				</div>
			{/each}
		</div>
	</div>
{/if}
