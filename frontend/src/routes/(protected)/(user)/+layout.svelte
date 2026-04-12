<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import InstallAppButton from '$lib/components/InstallAppButton.svelte';
	import LogoutIcon from '$lib/icons/LogoutIcon.svelte';
	import SettingsIcon from '$lib/icons/SettingsIcon.svelte';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	import LibraryIcon from '$lib/icons/LibraryIcon.svelte';
	import { page } from '$app/state';
	import ShowIcon from '$lib/icons/ShowIcon.svelte';
	import MovieIcon from '$lib/icons/MovieIcon.svelte';
	import CollectionIcon from '$lib/icons/CollectionIcon.svelte';

	let { children } = $props();

	const homeHref = resolve('/(protected)/(user)');
	const showsHref = resolve('/(protected)/(user)/shows');
	const moviesHref = resolve('/(protected)/(user)/movies');
	const collectionsHref = resolve('/(protected)/(user)/collections');

	const isActive = (href: string, includeChildren = false) => {
		if (page.url.pathname === href) return true;
		if (!includeChildren) return false;
		return page.url.pathname.startsWith(`${href}/`);
	};
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
<section id="selections" class="hidden items-center justify-center gap-4 p-4 md:flex">
	<div class="flex gap-2 rounded-full bg-neutral-800 px-3 py-2">
		<a
			href={homeHref}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 transition-all duration-100 md:text-lg"
			class:bg-neutral-700={isActive(homeHref)}
		>
			<HomeIcon />
			Home
		</a>
		<a
			href={showsHref}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 transition-all duration-100 md:text-lg"
			class:bg-neutral-700={isActive(showsHref, true)}
		>
			<ShowIcon />
			Shows
		</a>
		<a
			href={moviesHref}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 transition-all duration-100 md:text-lg"
			class:bg-neutral-700={isActive(moviesHref, true)}
		>
			<MovieIcon />
			Movies
		</a>
		<a
			href={collectionsHref}
			class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 transition-all duration-100 md:text-lg"
			class:bg-neutral-700={isActive(collectionsHref, true)}
		>
			<CollectionIcon />
			Collections
		</a>
	</div>
</section>
<section id="content" class="p-4 pb-24 md:pb-4">
	{@render children()}
</section>

<nav
	class="fixed right-0 bottom-0 left-0 border-t border-neutral-700 bg-neutral-900/95 px-2 pt-2 backdrop-blur-sm md:hidden"
	style="padding-bottom: max(0.5rem, env(safe-area-inset-bottom));"
>
	<div class="grid grid-cols-4 gap-1">
		<a
			href={homeHref}
			class="flex min-h-14 flex-col items-center justify-center gap-1 rounded-lg px-2 text-xs"
			class:bg-neutral-800={isActive(homeHref)}
			class:text-brand={isActive(homeHref)}
		>
			<HomeIcon />
			<span>Home</span>
		</a>
		<a
			href={showsHref}
			class="flex min-h-14 flex-col items-center justify-center gap-1 rounded-lg px-2 text-xs"
			class:bg-neutral-800={isActive(showsHref, true)}
			class:text-brand={isActive(showsHref, true)}
		>
			<LibraryIcon />
			<span>Shows</span>
		</a>
		<a
			href={moviesHref}
			class="flex min-h-14 flex-col items-center justify-center gap-1 rounded-lg px-2 text-xs"
			class:bg-neutral-800={isActive(moviesHref, true)}
			class:text-brand={isActive(moviesHref, true)}
		>
			<LibraryIcon />
			<span>Movies</span>
		</a>
		<a
			href={collectionsHref}
			class="flex min-h-14 flex-col items-center justify-center gap-1 rounded-lg px-2 text-xs"
			class:bg-neutral-800={isActive(collectionsHref, true)}
			class:text-brand={isActive(collectionsHref, true)}
		>
			<LibraryIcon />
			<span>Collections</span>
		</a>
	</div>
</nav>
