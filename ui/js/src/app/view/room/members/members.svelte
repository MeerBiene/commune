<script>
import { store } from '../../../store/store.js'
export let room;
import MemberItem from './member-item.svelte'

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

$: members = dm ? 
    account?.direct_messages?.filter(x => x.room_id == room.room_id)[0]?.members :
    $store.members[room?.server_id]

$: sorted = members ? Object.keys(members).sort().reduce(
  (obj, key) => { 
    obj[key] = members[key]; 
    return obj;
  }, 
  {}
) : null


$: chatroom = room?.room_type == 'chat'
$: topics = room?.room_type == 'topics'
$: thread = room?.room_type == 'thread'
$: dm = room?.room_type == 'dm'

</script>

<div class="members-c">
<div class="room-members pa2 scrl">
    {#if sorted}
        {#each Object.entries(sorted) as [user_id, info] (user_id)}
            <MemberItem info={info} user_id={user_id} />
        {/each}
    {/if}
</div>
</div>

<style>
.members-c {
    overflow: hidden;
    display: grid;
    grid-column: 2/3;
    grid-row: 2;
}
.room-members {
    background-color: var(--background-2);
    overflow: hidden auto;
}
.scrl  {
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: transparent transparent;
    scroll-behavior: smooth;
}

.scrl:hover {
    scrollbar-color: var(--background-1) transparent;
}

.scrl::-webkit-scrollbar {
  width: 4px;
    border-radius: 1px;
}
.scrl::-webkit-scrollbar-track {
    background: transparent;
}
.scrl::-webkit-scrollbar-thumb {
    background-color: transparent;
}
.scrl:hover::-webkit-scrollbar-thumb {
  background-color: var(--background-1);
}

</style>
