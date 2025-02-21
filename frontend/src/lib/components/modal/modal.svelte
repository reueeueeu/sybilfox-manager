<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { 
        children: Snippet;
    } = $props();
    let isActive = $derived(children.length != 0)
    console.log(children.length)
</script>

<div
	class="{isActive
		? 'absolute'
		: 'hidden'} top-0 z-40 flex h-full w-full flex-1 items-center bg-black/70 transition-all"
>
	{@render children()}
</div>
