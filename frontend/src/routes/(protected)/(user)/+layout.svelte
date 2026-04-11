<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import InstallAppButton from '$lib/components/InstallAppButton.svelte';
	import LogoutIcon from '$lib/icons/LogoutIcon.svelte';
	import SettingsIcon from '$lib/icons/SettingsIcon.svelte';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	import LibraryIcon from '$lib/icons/LibraryIcon.svelte';
	import { page } from '$app/state';

	let { children } = $props();

	const homeHref = resolve('/(protected)/(user)');
	const showsHref = resolve('/(protected)/(user)/shows');
	const moviesHref = resolve('/(protected)/(user)/movies');

	const isActive = (href: string) => page.url.pathname === href;
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
			href={homeHref}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 transition-all duration-100 md:text-lg"
			class:bg-neutral-700={isActive(homeHref)}
		>
			<HomeIcon /> Home
		</a>
		<a
			href={showsHref}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 transition-all duration-100 md:text-lg"
			class:bg-neutral-700={isActive(showsHref)}
		>
			<LibraryIcon />
			Shows
		</a>
		<a
			href={moviesHref}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 transition-all duration-100 md:text-lg"
			class:bg-neutral-700={isActive(moviesHref)}
		>
			<LibraryIcon />
			Movies
		</a>
		<a
			href={moviesHref}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 transition-all duration-100 md:text-lg"
			class:bg-neutral-700={isActive(moviesHref)}
		>
			<LibraryIcon />
			Collections
		</a>
	</div>
</section>
<section id="content" class="p-4">
	{@render children()}
</section>
