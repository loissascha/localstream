<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Home' }: Props = $props();
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
	<path d="M11.47 3.84a.75.75 0 0 1 1.06 0l7.5 7.5a.75.75 0 0 1-.53 1.28h-.75v7.13a1.5 1.5 0 0 1-1.5 1.5h-2.5a.75.75 0 0 1-.75-.75V16a.5.5 0 0 0-.5-.5h-3a.5.5 0 0 0-.5.5v4.5a.75.75 0 0 1-.75.75h-2.5a1.5 1.5 0 0 1-1.5-1.5v-7.13H4.5a.75.75 0 0 1-.53-1.28l7.5-7.5Z" />
</svg>
