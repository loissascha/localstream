<script lang="ts">
	import LastWatched from '$lib/components/LastWatched.svelte';
	import LastWatchedMovies from '$lib/components/LastWatchedMovies.svelte';
	import ShowListItem from '$lib/components/ShowListItem.svelte';
	import ItemGrid from '$lib/components/ItemGrid.svelte';
	import MovieListItem from '$lib/components/MovieListItem.svelte';
	import ShowIcon from '$lib/icons/ShowIcon.svelte';
	import MovieIcon from '$lib/icons/MovieIcon.svelte';
	import { loadMoviesDatabase, movies } from '$lib/movies.svelte';
	import { loadShowsDatabase, shows } from '$lib/shows.svelte';
	import { auth } from '$lib/auth.svelte';
	import ItemCarousel from '$lib/components/ItemCarousel.svelte';
	import ItemCarouselItem from '$lib/components/ItemCarouselItem.svelte';

	let latestMovies = $derived.by(() => {
		return [...movies.movies]
			.sort((a, b) => {
				const adate = new Date(a.created_at);
				const bdate = new Date(b.created_at);
				if (adate < bdate) return 1;
				if (adate > bdate) return -1;
				return 0;
			})
			.slice(0, 10);
	});
	let latestShows = $derived.by(() => {
		return [...shows.shows]
			.sort((a, b) => {
				const adate = new Date(a.created_at);
				const bdate = new Date(b.created_at);
				if (adate < bdate) return 1;
				if (adate > bdate) return -1;
				return 0;
			})
			.slice(0, 10);
	});

	$effect(() => {
		if (!auth.initialized) return;
		loadShowsDatabase().catch((e) => {
			const m = (e as Error).message;
			alert(m);
		});

		loadMoviesDatabase().catch((e) => {
			const m = (e as Error).message;
			alert(m);
		});
	});
</script>

<main>
	<LastWatched />
	<LastWatchedMovies />

	<section class="my-8">
		<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
			<ShowIcon />
			Recent Shows
		</h2>
		<ItemCarousel>
			{#each latestShows as show (show.id)}
				<ItemCarouselItem>
					<ShowListItem {show} />
				</ItemCarouselItem>
			{/each}
		</ItemCarousel>
	</section>

	<section class="my-8">
		<h2 class="mb-2 flex items-center gap-1 text-xl tracking-wider">
			<MovieIcon />
			Recent Movies
		</h2>
		<ItemCarousel>
			{#each latestMovies as movie (movie.id)}
				<ItemCarouselItem>
					<MovieListItem {movie} />
				</ItemCarouselItem>
			{/each}
		</ItemCarousel>
	</section>
</main>
