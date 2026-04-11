<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import InstallAppButton from '$lib/components/InstallAppButton.svelte';
	import LogoutIcon from '$lib/icons/LogoutIcon.svelte';
	import SettingsIcon from '$lib/icons/SettingsIcon.svelte';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	import LibraryIcon from '$lib/icons/LibraryIcon.svelte';

	let { children } = $props();
</script>

<section id="header" class="flex items-center justify-between px-4 py-4">
	<div class="flex items-center gap-2">
		<a href={resolve('/(protected)/(user)')} class="text-2xl font-semibold tracking-wider"
			>Localstream</a
		>
	</div>
	<div class="flex items-center gap-2">
		<InstallAppButton />
		{#if auth.isAdmin}
			<a
				href={resolve('/(protected)/(admin)/admin')}
				type="submit"
				class="cursor-pointer px-3 py-1.5 text-sm hover:text-brand"
			>
				<SettingsIcon />
			</a>
		{/if}
		<a
			href={resolve('/logout')}
			type="submit"
			class="cursor-pointer px-3 py-1.5 text-sm hover:text-brand"
		>
			<LogoutIcon />
		</a>
	</div>
</section>
<section id="seletions" class="flex items-center justify-center gap-4 p-4">
	<div class="flex gap-2 rounded-full bg-neutral-800 px-3 py-2">
		<a
			href={resolve('/(protected)/(user)')}
			class="flex cursor-pointer items-center gap-1 rounded-full bg-neutral-700 px-4 py-2 text-lg hover:bg-neutral-700"
		>
			<HomeIcon /> Home
		</a>
		<a
			href={resolve('/(protected)/(user)/shows')}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 text-lg hover:bg-neutral-700"
		>
			<LibraryIcon />
			Shows
		</a>
		<a
			href={resolve('/(protected)/(user)/movies')}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 text-lg hover:bg-neutral-700"
		>
			<LibraryIcon />
			Movies
		</a>
	</div>
</section>
<section id="content" class="p-4">
	{@render children()}
</section>
