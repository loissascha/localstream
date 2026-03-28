<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Sort' }: Props = $props();
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
		d="M8.22 6.22a.75.75 0 0 1 1.06 0l2.5 2.5a.75.75 0 0 1-1.06 1.06L9.5 8.56V18a.75.75 0 0 1-1.5 0V8.56l-1.22 1.22a.75.75 0 1 1-1.06-1.06l2.5-2.5Zm7.56 11.56a.75.75 0 0 1-1.06 0l-2.5-2.5a.75.75 0 0 1 1.06-1.06l1.22 1.22V6a.75.75 0 0 1 1.5 0v9.44l1.22-1.22a.75.75 0 1 1 1.06 1.06l-2.5 2.5Z"
	/>
</svg>
