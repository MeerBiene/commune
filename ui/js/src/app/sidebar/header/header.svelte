<script>
import { Route,Router } from 'svelte-navigator'
import Search from './search/search.svelte'
import ServerHeader from './server-header/server-header.svelte'

let matchSearch = `/` || `/messages`

</script>

<div class="sidebar-header">

    <Route path="/explore" primary={false} let:params>
        <div class="gr-default h-100">
            <div class="gr-center w-100 ph3 heading f3 bold">
                Explore
            </div>
        </div>
    </Route>

    <Route path="/:server/:room" primary={false} let:params>
        <ServerHeader server={params.server} room={params.room}/>
    </Route>

    <Route path="/:server/:room/message/:event" primary={false} let:params>
        <ServerHeader server={params.server} event={params.event} room={params.room}/>
    </Route>


    <Route path="/:server/:room/topic/:thread" primary={false} let:params>
        <ServerHeader server={params.server} room={params.room} thread={params.thread}/>
    </Route>

    <Route path="/messages/:sender" primary={false} let:params>
        <Search/>
    </Route>

    <Route path="/" primary={false} let:params>
        <Search/>
    </Route>
</div>


<style>
.sidebar-header {
    border-bottom: 1px solid var(--background-1);
}
.heading {
    color: var(--text);
}
</style>
