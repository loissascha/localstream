<script lang="ts">
	import { resolve } from '$app/paths';
	import type { MovieInfo } from '$lib/types/export_types';
	import ListItemA from './ListItemA.svelte';
	import MovieInfoDisplay from './MovieInfoDisplay.svelte';

	let { movie }: { movie: MovieInfo } = $props();
</script>

<ListItemA href={resolve('/(protected)/watch/movies/[movieID]', { movieID: movie.id })}>
	<MovieInfoDisplay {movie} />
</ListItemA>
