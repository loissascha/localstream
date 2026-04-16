<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Search' }: Props = $props();
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
	<path
		d="M10.5 3a7.5 7.5 0 1 0 4.76 13.29l4.22 4.22a.75.75 0 1 0 1.06-1.06l-4.22-4.22A7.5 7.5 0 0 0 10.5 3Z"
	/>
</svg>
