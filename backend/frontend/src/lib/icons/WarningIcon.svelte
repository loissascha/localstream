<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Warning' }: Props = $props();
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
		fill-rule="evenodd"
		clip-rule="evenodd"
		d="M10.7 4.08a1.5 1.5 0 0 1 2.6 0l7.26 12.9a1.5 1.5 0 0 1-1.3 2.24H4.74a1.5 1.5 0 0 1-1.3-2.24l7.26-12.9ZM12 8.5a1 1 0 0 1 1 1v4a1 1 0 1 1-2 0v-4a1 1 0 0 1 1-1Zm0 8a1.125 1.125 0 1 0 0-2.25 1.125 1.125 0 0 0 0 2.25Z"
	/>
</svg>
