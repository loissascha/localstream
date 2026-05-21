<script lang="ts">
	import { downloadMovieSubtitle, searchMovieSubtitles } from '$lib/api/movie_subtitles';
	import { auth } from '$lib/auth.svelte';
	import DownloadIcon from '$lib/icons/DownloadIcon.svelte';
	import SearchIcon from '$lib/icons/SearchIcon.svelte';
	import type { SubtitleProviderResult, MovieInfo } from '$lib/types/export_types';
	import Overlay from './Overlay.svelte';

	interface Props {
		close: () => void;
		movie: MovieInfo;
	}
	let { close, movie }: Props = $props();

	const INITIAL_LANGUAGE = 'EN';

	let searchingMetadata = $state(false);
	let searchQuery = $state('');
	let searchLang = $state(INITIAL_LANGUAGE);
	let subtitleResult = $state<SubtitleProviderResult[]>([]);

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.token) return;
		movie;
		searchMovieSubtitles(auth.token, movie.name, INITIAL_LANGUAGE)
			.then((result) => {
				subtitleResult = result;
			})
			.catch((e) => {
				const m = (e as Error).message;
				alert(m);
			});
	});

	function submitSearchForm() {
		try {
			if (searchingMetadata) return;
			if (!auth.token) return;
			searchingMetadata = true;
			if (searchQuery == '') {
				alert('No query');
				return;
			}
			searchMovieSubtitles(auth.token, searchQuery, searchLang)
				.then((result) => {
					subtitleResult = result;
				})
				.catch((e) => {
					const m = (e as Error).message;
					alert(m);
				})
				.finally(() => {
					searchingMetadata = false;
				});
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		}
	}

	async function downloadSubtitle(item: SubtitleProviderResult) {
		try {
			if (!auth.token) return;
			await downloadMovieSubtitle(auth.token, movie.id, item);
			close();
		} catch (e) {
			const m = (e as Error).message;
			alert(m);
		}
	}
</script>

<Overlay {close}>
	<h1 class="text-2xl font-bold tracking-wide">Subtitles {movie.name}</h1>
	<form
		onsubmit={(e) => {
			e.preventDefault();
			submitSearchForm();
		}}
		class="my-4 flex items-center gap-2"
	>
		<input
			bind:value={searchQuery}
			type="text"
			class="my-4 w-full rounded bg-neutral-700 px-4 py-2"
			placeholder="Search for new Metadata"
		/>
		<input bind:value={searchLang} type="text" class="my-4 w-30 rounded bg-neutral-700 px-4 py-2" />
		<button class="flex cursor-pointer gap-2 rounded bg-neutral-700 px-4 py-2 hover:bg-neutral-600">
			<SearchIcon /> Search
		</button>
	</form>
	{#each subtitleResult as subtitle}
		<div
			class="mb-1 flex items-center gap-1 border-b border-b-neutral-700 pb-1 last-of-type:border-b-0"
		>
			<div class="grow">
				<div>{subtitle.name} | {subtitle.name2}</div>
				<div class="font-serif text-sm">{subtitle.lang} ({subtitle.lang_short})</div>
				<div class="text-sm text-neutral-400">{subtitle.url}</div>
			</div>
			<div>
				<button
					class="cursor-pointer"
					onclick={() => {
						downloadSubtitle(subtitle);
					}}
				>
					<DownloadIcon />
				</button>
			</div>
		</div>
	{/each}
</Overlay>
