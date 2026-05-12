<script lang="ts">
	import type { MovieInfo } from '$lib/types/export_types';
	import Overlay from './Overlay.svelte';

	interface Props {
		close: () => void;
		movie: MovieInfo;
	}
	let { close, movie }: Props = $props();
</script>

<Overlay {close}>
	<h1 class="text-2xl font-bold tracking-wide">Subtitles {movie.name}</h1>
</Overlay>
