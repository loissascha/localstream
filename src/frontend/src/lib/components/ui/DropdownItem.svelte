<script lang="ts">
	let { children, href } = $props();
</script>

<a
	{href}
	role="menuitem"
	class="block px-4 py-2 text-left whitespace-nowrap transition-colors hover:bg-neutral-700 focus:bg-neutral-700 focus:outline-none"
>
	{@render children()}
</a>
