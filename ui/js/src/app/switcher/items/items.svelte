<script>
import { useLocation } from 'svelte-navigator'
import SwitcherItem from './item.svelte'
import DMItem from './dm-item.svelte'
import { store } from '../../store/store.js'

const location =  useLocation()

export let server;
export let room;

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

$: servers = account?.servers

$: expanded = $store.settings.switcher.mode == 'expanded'

$: ownedServers = servers.filter(server => server.owner)
$: joinedServers = servers.filter(server => !server.owner)

$: DMview = $location.pathname.includes('messages/')

$: dmNotifications = Object.entries(account?.dm_notifications)

</script>


<div class="items-c lh-copy">
    <div class="items pv1" 
        class:no-scrollbar={!expanded}
        class:scrl={expanded}>


        {#if expanded && dmNotifications.length > 0}
            <div class="spc-l">
                Direct Messages
            </div>
        {/if}

        {#each dmNotifications as [room_id, item] (room_id)}
            <DMItem 
            room_id={room_id}
            count={item.count} />
        {/each}

        {#if expanded && ownedServers.length > 0}
            <div class="spc-l">
                Your Servers
            </div>
        {/if}

        {#each ownedServers as s, i}
            <SwitcherItem 
            server={server}
            room={room}
            activeServer={s} />
        {/each}

        {#if expanded && joinedServers.length > 0}
            <div class="spc-l">
                Joined Servers
            </div>
        {/if}

        {#each joinedServers as s, i}
            <SwitcherItem 
            server={server}
            room={room}
            activeServer={s} />
        {/each}

        {#if expanded && servers.length == 0}
            <div class="noservers pt3">
                You haven't created or joined any servers yet.
            </div>
        {/if}

</div>
</div>

<style>
.items-c {
    overflow: hidden;
    display:grid;
}

.items {
    overflow: hidden scroll;
}

.scrl  {
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: transparent transparent;
    scroll-behavior: smooth;
}

.scrl:hover {
    scrollbar-color: var(--background-4) transparent;
}

.scrl::-webkit-scrollbar {
  width: 3px;
    margin-right: 2px;
    border-radius: 1px;
}
.scrl::-webkit-scrollbar-track {
    background: transparent;
}
.scrl::-webkit-scrollbar-thumb {
    background-color: transparent;
}

.scrl:hover::-webkit-scrollbar-thumb {
  background-color: var(--background-4);
}


.spc-l {
    text-transform: uppercase;
    font-size: 0.72rem;
    letter-spacing: 1px;
    font-weight: bold;
    color: var(--text-light);
    padding: 0.5rem 14px 0.25rem 14px;
}
.noservers {
    font-size: 0.9rem;
    font-weight: bold;
    letter-spacing: 1px;
    color: var(--text-muted);
    padding: 14px;
}
</style>
