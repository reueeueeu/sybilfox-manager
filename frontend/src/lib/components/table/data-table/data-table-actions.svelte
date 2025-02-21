<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import Play from 'lucide-svelte/icons/play';
	import Trash from 'lucide-svelte/icons/trash';
	import Pen from 'lucide-svelte/icons/settings';
	let { 
		id = 0,
		onRun,
		onEdit,
		onDelete
	}: {
		id: number, 
		onRun: ((id:number)=>void),
		onDelete: ((id:number)=>void),
		onEdit: ((id:number)=>void),
	} = $props();
</script>

<div class="ml-0 flex">
	<div class="flex">
		<Button variant="ghost" class="relative h-8 w-8" onclick={()=>onRun(id)}>
			<span class="sr-only">Start</span>
			<Play class="h-4 w-4"></Play>
		</Button>
		<Button variant="ghost" class="relative h-8 w-8" onclick={()=>onEdit(id)}>
			<span class="sr-only">Edit</span>
			<Pen class="h-4 w-4"></Pen>
		</Button>
		<Button variant="ghost" class=" relative h-8 w-8" onclick={()=>onDelete(id)}>
			<span class="sr-only">Delete</span>
			<Trash class="h-4 w-4"></Trash>
		</Button>
	</div>
</div>
