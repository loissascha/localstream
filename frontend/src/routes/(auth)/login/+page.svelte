<script lang="ts">
	import { resolve } from '$app/paths';
	import { API_URL } from '$lib/consts';
	import { onMount } from 'svelte';

	let loading = true;
	let data: any = null;
	let error: string | null = null;

	onMount(async () => {
		try {
			const res = await fetch(API_URL + '/auth/users/list');

			if (!res.ok) {
				throw new Error('Request failed');
			}

			data = await res.json();
			console.log('data:', data);
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	});
</script>

<main
	class="grid min-h-dvh place-items-center bg-[radial-gradient(circle_at_20%_10%,#d7e8ff_0%,transparent_40%),radial-gradient(circle_at_80%_0%,#c8f5e9_0%,transparent_32%),linear-gradient(180deg,#eef2f7_0%,#dce6f2_100%)] px-4 py-6"
>
	<section
		class="w-full max-w-sm rounded-2xl border border-slate-900/10 bg-white/85 p-6 shadow-lg shadow-slate-900/10"
	>
		<h1 class="m-0 text-2xl font-semibold text-slate-900">Choose Profile</h1>
		{#if loading}
			<p>Loading</p>
		{:else if error}
			<p>Error: {error}</p>
		{:else}
			<div class="my-8 flex gap-2 items-center justify-center">
				{#each data as item}
					<div
						class="h-22 w-22 cursor-pointer place-content-center place-items-center rounded border border-neutral-300 bg-neutral-200 text-center shadow"
					>
						{item.username}
					</div>
				{/each}
			</div>
		{/if}

		<p class="text-sm">
			No profile yet? <a class="cursor-pointer text-blue-600" href={resolve('/(auth)/register')}
				>create a new one</a
			>
		</p>
	</section>
</main>
