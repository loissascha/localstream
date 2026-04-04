<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { listLatestWatchstateByMovie } from '$lib/api/watchstate';
	import { auth } from '$lib/auth.svelte';
	import ChevronRightIcon from '$lib/icons/ChevronRightIcon.svelte';
	import { type WatchstateMovieResponse } from '$lib/types/export_types';

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
	<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
		<ChevronRightIcon />
		Continue Movies
	</h2>
	<div class="flex gap-4">
		{#each data as d}
			<a
				href={resolve('/(protected)/(watch)/movies/[movieID]', {
					movieID: d.movie_id
				})}
				class="flex w-60 cursor-pointer flex-col justify-between gap-2 rounded bg-neutral-800 p-4 hover:bg-neutral-700"
			>
				<div class="font-bold">{d.movie_info.name}</div>
				<div>
					<div class="bg-neutral-600">
						<div style={`width: ${d.percentage}%;`} class={`h-2 bg-blue-300 text-sm`}></div>
					</div>
				</div>
			</a>
		{/each}
	</div>
</section>
