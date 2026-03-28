<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';

	let { children } = $props();
</script>

<section id="header" class="flex items-center justify-between px-4 py-4">
	<div>
		<a href={resolve('/(protected)/(user)')} class="text-2xl font-semibold tracking-wider"
			>Localstream</a
		>
	</div>
	<div class="flex items-center gap-2">
		{#if auth.isAdmin}
			<a
				href={resolve('/(protected)/(admin)/admin')}
				type="submit"
				class="cursor-pointer rounded-md border border-neutral-500 bg-neutral-600 px-3 py-1.5 text-sm hover:bg-neutral-500"
			>
				Admin
			</a>
		{/if}
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
