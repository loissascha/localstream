<script lang="ts">
	import { resolve } from '$app/paths';
	import { auth } from '$lib/auth.svelte';
	import { searchLibrary } from '$lib/api/search';
	import InstallAppButton from '$lib/components/InstallAppButton.svelte';
	import HomeIcon from '$lib/icons/HomeIcon.svelte';
	import { page } from '$app/state';
	import ShowIcon from '$lib/icons/ShowIcon.svelte';
	import MovieIcon from '$lib/icons/MovieIcon.svelte';
	import CollectionIcon from '$lib/icons/CollectionIcon.svelte';
	import SearchIcon from '$lib/icons/SearchIcon.svelte';
	import type { SearchResponse } from '$lib/types/export_types';
	import Logo from '$lib/components/ui/Logo.svelte';
	import LayoutUserDropdown from '$lib/components/ui/LayoutUserDropdown.svelte';

	let { children } = $props();
	let searchRoot: HTMLElement | null = $state(null);
	let searchQuery = $state('');
	let searchOpen = $state(false);
	let searchLoading = $state(false);
	let searchError = $state('');
	let searchResults = $state<SearchResponse>({ shows: [], movies: [] });

	let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null;
	let latestSearchRequest = 0;

	const homeHref = resolve('/(protected)/(user)');
	const showsHref = resolve('/(protected)/(user)/shows');
	const moviesHref = resolve('/(protected)/(user)/movies');
	const collectionsHref = resolve('/(protected)/(user)/collections');
	const hasSearchResults = $derived(
		searchResults.shows.length > 0 || searchResults.movies.length > 0
	);
	const shouldShowSearchDropdown = $derived(
		searchOpen &&
			(searchLoading || searchError !== '' || hasSearchResults || searchQuery.trim().length >= 3)
	);

	const isActive = (href: string, includeChildren = false) => {
		if (page.url.pathname === href) return true;
		if (!includeChildren) return false;
		return page.url.pathname.startsWith(`${href}/`);
	};

	function resetSearchResults() {
		searchResults = { shows: [], movies: [] };
		searchError = '';
		searchLoading = false;
	}

	function closeSearch() {
		searchOpen = false;
	}

	function openSearch() {
		searchOpen = true;
	}

	async function runSearch(query: string) {
		if (!auth.token) return;

		const requestId = ++latestSearchRequest;
		searchLoading = true;
		searchError = '';

		try {
			const result = await searchLibrary(auth.token, query);
			if (requestId !== latestSearchRequest) return;

			searchResults = result;
			searchOpen = true;
		} catch (error) {
			if (requestId !== latestSearchRequest) return;

			resetSearchResults();
			searchError = error instanceof Error ? error.message : 'Failed to search library';
			searchOpen = true;
		} finally {
			if (requestId === latestSearchRequest) {
				searchLoading = false;
			}
		}
	}

	function handleSearchInput() {
		if (searchDebounceTimer) {
			clearTimeout(searchDebounceTimer);
		}

		const query = searchQuery.trim();
		if (query.length < 3) {
			latestSearchRequest += 1;
			resetSearchResults();
			if (query.length === 0) {
				searchOpen = false;
			} else {
				searchOpen = true;
			}
			return;
		}

		searchDebounceTimer = setTimeout(() => {
			runSearch(query);
		}, 250);
	}

	function handleSearchSubmit(event: SubmitEvent) {
		event.preventDefault();
		const query = searchQuery.trim();
		if (query.length < 3) {
			return;
		}

		if (searchDebounceTimer) {
			clearTimeout(searchDebounceTimer);
		}

		runSearch(query);
	}

	function handleSearchKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			closeSearch();
		}
	}

	function handleSearchResultClick() {
		closeSearch();
	}

	$effect(() => {
		return () => {
			if (searchDebounceTimer) {
				clearTimeout(searchDebounceTimer);
			}
		};
	});

	$effect(() => {
		const handlePointerDown = (event: MouseEvent) => {
			if (!searchRoot) return;
			const target = event.target;
			if (!(target instanceof Node)) return;
			if (!searchRoot.contains(target)) {
				closeSearch();
			}
		};

		document.addEventListener('mousedown', handlePointerDown);

		return () => {
			document.removeEventListener('mousedown', handlePointerDown);
		};
	});

	$effect(() => {
		const currentPath = page.url.pathname;
		if (currentPath != null) {
			closeSearch();
		}
	});
</script>

<section
	id="header"
	class="sticky top-0 z-40 flex items-center justify-between bg-neutral-900 px-4 py-4 md:grid md:grid-cols-[1fr_auto_1fr]"
>
	<div class="flex grow items-center gap-2">
		<a
			href={resolve('/(protected)/(user)')}
			class="flex items-center gap-2 text-2xl font-semibold tracking-wider select-none"
			><Logo />ocalstream</a
		>
	</div>
	<section id="middle" class="flex items-center justify-center gap-2 px-4">
		<div class="hidden flex-col gap-2 lg:flex lg:flex-row">
			<a
				href={homeHref}
				class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 text-neutral-400 transition-all duration-100 md:text-lg"
				class:text-white={isActive(homeHref)}
			>
				<HomeIcon />
				Home
			</a>
			<a
				href={showsHref}
				class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 text-neutral-400 transition-all duration-100 md:text-lg"
				class:text-white={isActive(showsHref, true)}
			>
				<ShowIcon />
				Shows
			</a>
			<a
				href={moviesHref}
				class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 text-neutral-400 transition-all duration-100 md:text-lg"
				class:text-white={isActive(moviesHref, true)}
			>
				<MovieIcon />
				Movies
			</a>
			<a
				href={collectionsHref}
				class="flex cursor-pointer items-center gap-1 rounded-full px-4 py-2 text-neutral-400 transition-all duration-100 md:text-lg"
				class:text-white={isActive(collectionsHref, true)}
			>
				<CollectionIcon />
				Collections
			</a>
		</div>
	</section>
	<div class="flex grow items-center justify-end gap-2">
		<div class="relative w-full max-w-80" bind:this={searchRoot}>
			<form class="flex items-center gap-2" onsubmit={handleSearchSubmit}>
				<div class="relative flex-1">
					<input
						type="text"
						placeholder="Search shows and movies"
						class="w-full rounded-full border border-transparent bg-neutral-800 px-4 py-2 pr-11 transition outline-none focus:border-neutral-600"
						bind:value={searchQuery}
						oninput={handleSearchInput}
						onfocus={openSearch}
						onkeydown={handleSearchKeydown}
					/>
					<button
						type="submit"
						class="absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer text-neutral-300 hover:text-white"
						aria-label="Search"
					>
						<SearchIcon />
					</button>
				</div>
			</form>

			{#if shouldShowSearchDropdown}
				<div
					class="absolute top-full right-0 left-0 z-40 mt-2 overflow-hidden rounded-2xl border border-neutral-700 bg-neutral-900/98 shadow-2xl backdrop-blur-sm"
				>
					<div class="max-h-[70vh] overflow-y-auto p-2">
						{#if searchLoading}
							<div class="px-3 py-4 text-sm text-neutral-300">Searching...</div>
						{:else if searchError}
							<div class="px-3 py-4 text-sm text-red-400">{searchError}</div>
						{:else if searchQuery.trim().length < 3}
							<div class="px-3 py-4 text-sm text-neutral-400">Type at least 3 characters</div>
						{:else if !hasSearchResults}
							<div class="px-3 py-4 text-sm text-neutral-400">No results</div>
						{:else}
							{#if searchResults.shows.length > 0}
								<div class="pb-2">
									<div
										class="px-3 py-2 text-xs font-semibold tracking-[0.18em] text-neutral-500 uppercase"
									>
										Shows
									</div>
									<div class="space-y-1">
										{#each searchResults.shows as show (show.id)}
											<a
												href={resolve('/(protected)/(user)/shows/[showID]', { showID: show.id })}
												class="flex items-center justify-between gap-3 rounded-xl px-3 py-3 transition hover:bg-neutral-800"
												onclick={handleSearchResultClick}
											>
												<div class="flex min-w-0 items-center gap-3">
													<div class="rounded-lg bg-neutral-800 p-2 text-neutral-200">
														{#if show.medium_image_url != ''}
															<img alt={show.name} src={show.medium_image_url} class="w-10" />
														{:else}
															<ShowIcon />
														{/if}
													</div>
													<div class="min-w-0">
														<div class="truncate font-medium text-white">{show.name}</div>
														<div class="text-sm text-neutral-400">
															{show.year > 0 ? show.year : 'Show'}
														</div>
													</div>
												</div>
												<div class="text-xs tracking-[0.18em] text-neutral-500 uppercase">Show</div>
											</a>
										{/each}
									</div>
								</div>
							{/if}

							{#if searchResults.movies.length > 0}
								<div>
									<div
										class="px-3 py-2 text-xs font-semibold tracking-[0.18em] text-neutral-500 uppercase"
									>
										Movies
									</div>
									<div class="space-y-1">
										{#each searchResults.movies as movie (movie.id)}
											<a
												href={resolve('/(protected)/(user)/movies/[movieID]', {
													movieID: movie.id
												})}
												class="flex items-center justify-between gap-3 rounded-xl px-3 py-3 transition hover:bg-neutral-800"
												onclick={handleSearchResultClick}
											>
												<div class="flex min-w-0 items-center gap-3">
													<div class="rounded-lg bg-neutral-800 p-2 text-neutral-200">
														{#if movie.medium_image_url != ''}
															<img alt={movie.name} src={movie.medium_image_url} class="w-10" />
														{:else}
															<MovieIcon />
														{/if}
													</div>
													<div class="min-w-0">
														<div class="truncate font-medium text-white">{movie.name}</div>
														<div class="text-sm text-neutral-400">
															{movie.year > 0 ? movie.year : 'Movie'}
														</div>
													</div>
												</div>
												<div class="text-xs tracking-[0.18em] text-neutral-500 uppercase">
													Movie
												</div>
											</a>
										{/each}
									</div>
								</div>
							{/if}
						{/if}
					</div>
				</div>
			{/if}
		</div>
		<InstallAppButton />
		<LayoutUserDropdown />
	</div>
</section>
<section id="content" class="p-4 pb-24 md:pb-4">
	{@render children()}
</section>

<nav
	class="fixed right-0 bottom-0 left-0 border-t border-neutral-700 bg-neutral-900/95 px-2 pt-2 backdrop-blur-sm lg:hidden"
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
			<ShowIcon />
			<span>Shows</span>
		</a>
		<a
			href={moviesHref}
			class="flex min-h-14 flex-col items-center justify-center gap-1 rounded-lg px-2 text-xs"
			class:bg-neutral-800={isActive(moviesHref, true)}
			class:text-brand={isActive(moviesHref, true)}
		>
			<MovieIcon />
			<span>Movies</span>
		</a>
		<a
			href={collectionsHref}
			class="flex min-h-14 flex-col items-center justify-center gap-1 rounded-lg px-2 text-xs"
			class:bg-neutral-800={isActive(collectionsHref, true)}
			class:text-brand={isActive(collectionsHref, true)}
		>
			<CollectionIcon />
			<span>Collections</span>
		</a>
	</div>
</nav>
