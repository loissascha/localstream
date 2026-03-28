<script lang="ts">
	let { children } = $props();

	$effect(() => {});
</script>

<div>
	<div>Layout</div>
	<div>
		{@render children()}
	</div>
</div>
