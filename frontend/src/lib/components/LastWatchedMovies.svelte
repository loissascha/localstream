<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { listLatestWatchstateByMovie } from '$lib/api/watchstate_movie';
	import { auth } from '$lib/auth.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import { type WatchstateMovieResponse } from '$lib/types/export_types';
	import ItemGrid from './ItemGrid.svelte';
	import ListItemA from './ListItemA.svelte';
	import MovieInfoDisplay from './MovieInfoDisplay.svelte';

	let data = $state<WatchstateMovieResponse[]>([]);

	async function updateData() {
		if (!auth.token) {
			return;
		}
		try {
			const result = await listLatestWatchstateByMovie(auth.token);
			data = result;
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.loggedIn) {
			goto(resolve('/(auth)/login'));
			return;
		}
		updateData();
	});
</script>

<section class="my-4">
	{#if data.length > 0}
		<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
			<ChevronRightIcon />
			Continue Movies
		</h2>
		<ItemGrid>
			{#each data as d (d.id)}
				<ListItemA
					href={resolve('/(protected)/(watch)/movies/[movieID]', {
						movieID: d.movie_id
					})}
				>
					<div>
						<MovieInfoDisplay movie={d.movie_info} />
					</div>
					<div>
						<div class="bg-neutral-600">
							<div style={`width: ${d.percentage}%;`} class={`h-2 bg-brand text-sm`}></div>
						</div>
					</div>
				</ListItemA>
			{/each}
		</ItemGrid>
	{/if}
</section>
