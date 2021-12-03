<script>
import { onMount, tick } from 'svelte'
import { useLocation } from 'svelte-navigator'
import { thread, threadBig } from '../utils/icons'
import { store } from '../store/store.js'
import { closeBig  } from '../utils/icons.js'
import EventItem from '../view/room/events/event-item.svelte'
import Input from '../view/room/input/input.svelte'
import Room  from '../view/room/room.svelte'

export let room;

function kill() {
    store.closeThread(room)
}

const location = useLocation()
$: server = $location.pathname.split("/")[1]
$: page = {
    pathname: `/${server}${room?.thread?.thread_room?.pathname}`,
    visible: false,
}





</script>

<div class="header flex">
    <div class="gr-center ph3">
        {@html thread}
    </div>
    <div class="gr-center flex-one">
        <strong>{room?.thread?.thread_room?.name}</strong>
    </div>
    <div class="gr-center ph3 pointer icon" 
        on:click={kill}
        aria-label="Close"
        data-microtip-position="bottom"
        data-microtip-size="fit"
        role="tooltip">
        {@html closeBig}
    </div>
</div>

<div></div>


<Room threadView={true} page={page}/>

<Input threadView={true} page={page} />


<style>
.header {
    border-bottom: 1px solid var(--background-1);
}

</style>
