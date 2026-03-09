<script lang="ts">
	import { page } from '$app/state';

	const showId = $derived(page.params.showID ?? '');
	const episodeId = $derived(page.params.episodeID ?? '');
</script>

<main>
	Show: {showId} <br />
	Episode: {episodeId}
</main>
