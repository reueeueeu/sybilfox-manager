<script lang="ts">
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import Table from './table/table.svelte';
	import Button from './ui/button/button.svelte';
	import Modal from './modal/modal.svelte';
	import ProfileEditor from './modal/profile-editor.svelte';
	import { profile, proxy } from '$lib/wailsjs/go/models';
	import { tapOutside } from 'svelte-outside';
	import {
		HideConfig,
		InstallBrowser,
		IsBrowserInstalled,
		ListConfigs,
		Logout,
		NewConfig,
		RunConfig,
		UpdateConfig
	} from '$lib/wailsjs/go/main/App';
	import { onMount } from 'svelte';
	import LogOut from 'lucide-svelte/icons/log-out';
	import { toast } from 'svelte-sonner';
	import { Root } from './ui/avatar';
	import * as Card from './ui/card';

	//let configs = new SvelteSet<profile.Config>([]);
	let configs: profile.Config[] = $state([]);
	let menuSelected = $state('profiles');
	let prevSelected = $state('');
	$effect(() => {
		prevSelected = menuSelected;
	});
	function getProfileConfig(id: number) {
		return configs.find((c) => {
			return c.id == id;
		});
	}
	function capitalize(str: string) {
		return str.charAt(0).toUpperCase() + str.slice(1);
	}
	let installingBrowser = $state(false);
	let profileRequest: profile.Request | undefined = $state();
	let modalActive = $derived(profileRequest !== undefined || installingBrowser);
	$effect(() => {
		console.log(modalActive);
	});
	onMount(async () => {
		//(await ListConfigs()).map((c) => {configs.add(c)});
		configs = await ListConfigs();
		console.log(configs);
	});
	let { onLogout }: { onLogout: () => void } = $props();
	let profileId = -1;
</script>

<Sidebar.Provider>
	<Sidebar.Inset class="relative">
		<header
			class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12"
		>
			<div class="flex items-center gap-2 px-4">
				<Breadcrumb.Root>
					<Breadcrumb.List>
						<Breadcrumb.Item class="hidden md:block">
							<Breadcrumb.Link href="#">{capitalize(menuSelected)}</Breadcrumb.Link>
						</Breadcrumb.Item>
					</Breadcrumb.List>
				</Breadcrumb.Root>
			</div>
			<div class="z-50 ml-auto mr-4 flex gap-4">
				<Button
					onclick={() => {
						profileRequest = new profile.Request({
							name: '',
							access_code: -1,
							proxy: new proxy.Config({
								host: '',
								user: '',
								password: ''
							})
						});
					}}>Add Profile</Button
				>
				<Button
					variant="outline"
					onclick={() => {
						Logout().then(() => {
							onLogout();
						});
					}}><LogOut></LogOut></Button
				>
			</div>
		</header>
		<div class="absolute flex h-full w-full">
			<div class="m-auto">
				<img alt="logo" class="h-[256px] w-[256px]" src="/logo.png" />
				<!-- <p class="text-4xl font-extrabold">Sybilfox Manager</p> -->
			</div>
		</div>
		<Table
			{configs}
			onProfileDelete={async (id) => {
				try {
					await HideConfig(id);
					configs = configs.filter((config) => config.id !== id);
				} catch (e) {
					toast("Can't delete config. Error: " + e);
				}
			}}
			onProfileEdit={(id) => {
				let config = configs.find((c) => {
					return c.id == id;
				});
				profileRequest = config?.request;
				profileId = id;
			}}
			onProfileRun={async (id) => {
				let config = configs.find((c) => {
					return c.id == id;
				});
				console.log(config);
				if (config?.request.imported) {
					profileId = id;
					profileRequest = config.request;
					return;
				}
				if (!(await IsBrowserInstalled())) {
					installingBrowser = true;
					try {
						await InstallBrowser();
					} catch (e) {
						toast('Failed to download browser. Error: ' + e);
						return;
					} finally {
						installingBrowser = false;
						profileRequest = undefined;
					}
				}
				try {
					await RunConfig(id);
				} catch (e) {
					toast("Can't start browser. Error: " + e);
				} finally {
					installingBrowser = false;
					profileRequest = undefined;
				}
			}}
		/>
		{#if modalActive}
			<Modal>
				<div
					class="mx-auto"
					use:tapOutside={(e) => {
						profileRequest = undefined;
					}}
				>
					{#if installingBrowser}
						<div class="">
							<Card.Root><Card.Content>Installing Browser...</Card.Content></Card.Root>
						</div>
					{/if}
					{#if profileRequest !== undefined}
						<ProfileEditor
							request={profileRequest}
							onSubmit={async (request) => {
								async function update() {
									let updated = await UpdateConfig(profileId, request);
									const index = configs.findIndex((config) => config.id === updated.id);
									if (index !== -1) {
										// Replace the existing config
										configs = configs.map((config) =>
											config.id === updated.id ? updated : config
										);
									} else {
										// Add the updated config if it doesn't exist
										configs = [...configs, updated];
									}
								}
								try {
									console.log(request);
									if (request.imported) {
										await update();
										if (!(await IsBrowserInstalled())) {
											installingBrowser = true;
											try {
												await InstallBrowser();
											} catch (e) {
												toast('Failed to download browser. Error: ' + e);
												return;
											} finally {
												installingBrowser = false;
												profileRequest = undefined;
											}
										}
										try {
											await RunConfig(profileId);
										} catch (e) {
											toast("Can't start browser. Error: " + e);
											return;
										} finally {
											installingBrowser = false;
											profileRequest = undefined;
										}
										return;
									}
									if (request.access_code == -1) {
										let profile = await NewConfig(request.name, request.proxy);
										configs = [...configs, profile];
										console.log(profile);
									} else if (request.access_code > 0) {
										await update();
									}
									profileRequest = undefined;
								} catch (e) {
									toast("Can't create profile. Error: " + e);
								}
								return;
							}}
						></ProfileEditor>
					{/if}
				</div>
			</Modal>
		{/if}
	</Sidebar.Inset>
</Sidebar.Provider>
