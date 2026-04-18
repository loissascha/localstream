<script lang="ts">
	import { resolve } from '$app/paths';
	import type { MovieInfo } from '$lib/types/export_types';
	import ListItemA from './ListItemA.svelte';
	import MovieInfoDisplay from './MovieInfoDisplay.svelte';

	interface Props {
		movie: MovieInfo;
		selectable?: boolean;
		selected?: boolean;
	}

	let { movie, selectable = false, selected = $bindable(false) }: Props = $props();
</script>

<div class="relative">
	<ListItemA href={resolve('/(protected)/(user)/movies/[movieID]', { movieID: movie.id })}>
		<MovieInfoDisplay {movie} />
	</ListItemA>

	{#if selectable}
		<button
			type="button"
			class="absolute top-2 right-2 z-10"
			role="checkbox"
			aria-checked={selected}
			aria-label={`Select ${movie.name}`}
			onclick={(event: MouseEvent) => {
				event.preventDefault();
				event.stopPropagation();
				selected = !selected;
			}}
		>
			<span
				class={`flex h-7 w-7 items-center justify-center rounded-full border shadow-sm transition-all duration-150 ${selected ? 'border-brand bg-brand text-white' : 'border-neutral-500/80 bg-neutral-950/85 text-transparent hover:border-neutral-300 hover:bg-neutral-900'}`}
			>
				<svg aria-hidden="true" class="h-4 w-4" viewBox="0 0 16 16" fill="none">
					<path
						d="M3.5 8.5L6.5 11.5L12.5 4.5"
						stroke="currentColor"
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
					/>
				</svg>
			</span>
		</button>
	{/if}
</div>
