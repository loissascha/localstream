<script lang="ts">
	import { page } from '$app/state';

	const showId = $derived(page.params.showID ?? '');
	const seasonId = $derived(page.params.seasonID ?? '');
	const episodeId = $derived(page.params.episodeID ?? '');
</script>

<main>
	Show: {showId} <br />
	Season: {seasonId}<br />
	Episode: {episodeId}
</main>
