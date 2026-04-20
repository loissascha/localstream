<script lang="ts">
	let { children, href } = $props();
</script>

<a {href} class="px-4 py-2">
	{@render children()}
</a>
