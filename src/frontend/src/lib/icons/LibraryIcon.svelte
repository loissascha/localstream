<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Library' }: Props = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	viewBox="0 0 24 24"
	fill="currentColor"
	width={size}
	height={size}
	class={className}
	role="img"
	aria-hidden={title ? undefined : 'true'}
>
	{#if title}
		<title>{title}</title>
	{/if}
	<rect x="4" y="4" width="4" height="16" rx="1" />
	<rect x="10" y="4" width="4" height="16" rx="1" />
	<rect x="16" y="4" width="4" height="16" rx="1" />
</svg>
