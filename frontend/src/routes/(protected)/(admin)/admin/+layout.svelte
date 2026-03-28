<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let { children } = $props();

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.isAdmin) {
			goto(resolve('/(protected)'));
		}
	});
</script>

<section id="header" class="flex items-center justify-between bg-neutral-900 px-4 py-4">
	<div>
		<a href={resolve('/(protected)/(user)')} class="text-2xl font-semibold tracking-wider"
			>Localstream Admin</a
		>
	</div>
	<div>
		<a
			href={resolve('/logout')}
			type="submit"
			class="cursor-pointer rounded-md border border-neutral-500 bg-neutral-600 px-3 py-1.5 text-sm hover:bg-neutral-500"
		>
			Log out
		</a>
	</div>
</section>
<section id="content" class="p-4">
	{@render children()}
</section>
