<script lang="ts">
	import Button from "$lib/components/ui/button/button.svelte";

	import type { WithElementRef } from "bits-ui";
	import type { HTMLAnchorAttributes, HTMLButtonAttributes } from "svelte/elements";

    let {
        text = '',
        ...props
    }: WithElementRef<HTMLButtonAttributes> &
		WithElementRef<HTMLAnchorAttributes> & {text: string} = $props()
</script>
<a class="{props.class} cursor-pointer hover:text-primary" {...props}>{text}</a>