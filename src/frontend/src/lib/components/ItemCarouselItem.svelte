<script lang="ts">
	let { children } = $props();
</script>

<div class="w-40 shrink-0 md:w-50 lg:w-60">
	{@render children()}
</div>
