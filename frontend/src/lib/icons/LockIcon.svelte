<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Lock' }: Props = $props();
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
		d="M7.5 10V8a4.5 4.5 0 1 1 9 0v2h.5A2.5 2.5 0 0 1 19.5 12.5v7A2.5 2.5 0 0 1 17 22H7a2.5 2.5 0 0 1-2.5-2.5v-7A2.5 2.5 0 0 1 7 10h.5Zm1.5 0h6V8a3 3 0 1 0-6 0v2Z"
	/>
	<circle cx="12" cy="15.5" r="1.25" />
	<rect x="11.25" y="16.5" width="1.5" height="2.5" rx="0.75" />
</svg>
