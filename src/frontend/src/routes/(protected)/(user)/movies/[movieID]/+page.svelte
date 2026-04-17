<script lang="ts">
	import { page } from '$app/state';

	const movieId = $derived(page.params.movieID ?? '');
</script>

Movie
