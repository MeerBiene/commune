<script>
import { Router, Route } from 'svelte-navigator'
import Rooms from './rooms/rooms.svelte'
import Home from './home/home.svelte'


</script>

<Router primary={false}>
<div class="sidebar-view">

    <Route path="/explore" let:params>
    </Route>
    <Route path="/:server/:room" let:params>
        <Rooms 
        room={params.room}
        server={params.server}/>
    </Route>
    <Route path="/:server/:room/message/:event" let:params>
        <Rooms 
        room={params.room}
        event={params.event}
        server={params.server}/>
    </Route>
    <Route path="/:server/:room/topic/:thread" let:params>
        <Rooms 
        room={params.room}
        server={params.server}
        thread={params.thread}/>
    </Route>
    <Route path="/messages/:sender" let:params>
        <Home 
        sender={params.sender} />
    </Route>
    <Route path="/" let:params>
        <Home  />
    </Route>
</div>
</Router>

<style>
.sidebar-view {
    overflow: hidden;
    display: grid;
}
</style>
